package patch

import (
	"errors"
	"fmt"
	"git.code.oa.com/goom/mocker/internal/logger"
	"reflect"
	"strings"
	"sync"
)

// patch is an applied patch
// needed to undo a patch
type patch struct {
	originalBytes []byte
	targetPtr     uintptr
	replacement   *reflect.Value
	originPtr     uintptr
}

var (
	lock      = sync.Mutex{}
	patches   = make(map[uintptr]patch)
	ptrholder = make(map[uintptr]interface{})
	redirect  = make(map[uintptr]uintptr)
)

// PatchGuard 代理执行控制句柄, 可通过此对象进行代理还原
type PatchGuard struct {
	target      uintptr
	replacement reflect.Value
	originFunc  uintptr
	jumpData    []byte
	applied     bool
}

func PatchLock() {
	lock.Lock()
}

func PatchUnlock() {
	lock.Unlock()
}

func (g *PatchGuard) Apply() {
	g.applied = true
	// 执行函数调用地址替换(延迟执行)
	if err := copyToLocation(g.target, g.jumpData); err != nil {
		logger.LogWarningf("Apply to 0x%x error: %s", g.target, err)
	}
	ShowInst(fmt.Sprintf("apply copy to 0x%x", g.target), g.target, 20, logger.DebugLevel)
}

// Unpatch 取消代理,还原指令码
// 外部调用请使用PatchGuard.UnpatchWithLock()
func (g *PatchGuard) Unpatch() {
	if g != nil && g.applied {
		unpatchValue(g.target)
	}
}

// UnpatchWithLock 外部调用需要加锁
func (g *PatchGuard) UnpatchWithLock() {
	PatchLock()
	defer PatchUnlock()

	if g != nil && g.applied {
		unpatchValue(g.target)
	}
}

// Restore 重新应用代理
func (g *PatchGuard) Restore() {
	if g != nil && g.applied {
		_, _ = PatchPtr(g.target, g.replacement)
	}
}

// OriginFunc 获取应用代理后的原函数地址(和代理前的原函数地址不一样)
func (g *PatchGuard) OriginFunc() uintptr {
	return g.originFunc
}

// Patch 将函数调用指定代理函数
// target 原始函数
// replacement 代理函数
func Patch(target, replacement interface{}) (*PatchGuard, error) {
	return PatchTrampoline(target, replacement, nil)
}

// Patch 将函数调用指定代理函数
// target 原始函数
// replacement 代理函数
// trampoline 指定跳板函数(可不指定,传nil)
func PatchTrampoline(target, replacement interface{}, trampoline interface{}) (*PatchGuard, error) {
	t := reflect.ValueOf(target)
	r := reflect.ValueOf(replacement)
	originFunc, jumpData, err := patchValue(t, r, trampoline)
	if err != nil {
		return nil, err
	}
	ptrholder[t.Pointer()] = replacement

	return &PatchGuard{t.Pointer(), r, originFunc, jumpData, false}, nil
}

func UnsafePatch(target, replacement interface{}) (*PatchGuard, error) {
	return UnsafePatchTrampoline(target, replacement, nil)
}

// Patch replaces a function with another
func UnsafePatchTrampoline(target, replacement interface{}, trampoline interface{}) (*PatchGuard, error) {
	t := reflect.ValueOf(target)
	r := reflect.ValueOf(replacement)

	trampolinePtr, err := getTrampolinePtr(trampoline)
	if err != nil {
		return nil, err
	}

	originFunc, jumpData, err := unsafePatchValue(t, r, trampolinePtr)
	if err != nil {
		return nil, err
	}
	ptrholder[t.Pointer()] = replacement

	return &PatchGuard{t.Pointer(), r, originFunc, jumpData, false}, nil
}

// patchValue 对value进行应用代理
func patchValue(target, replacement reflect.Value, trampoline interface{}) (uintptr, []byte, error) {

	// 外部手动check(使用signature包) modified by @jake
	//if target.Type() != replacement.Type() {
	//	return 0, errors.New(fmt.Sprintf("target and replacement have to have the same type %s != %s", target.Type(), replacement.Type()))
	//}

	// 外部手动check(使用signature包) modified by @jake
	//if target.Type() != reflect.TypeOf(trampoline).Elem() {
	//	return 0, errors.New(fmt.Sprintf("target and trampoline have to have the same type %s != %s", target.Type(), trampolineElemType))
	//}

	trampolinePtr, err := getTrampolinePtr(trampoline)
	if err != nil {
		return 0, nil, err
	}

	ptrholder[trampolinePtr] = trampoline
	return unsafePatchValue(target, replacement, trampolinePtr)
}

// unsafePatchValue 不做类型检查
func unsafePatchValue(target, replacement reflect.Value, trampoline uintptr) (uintptr, []byte, error) {
	if target.Kind() != reflect.Func {
		return 0, nil, errors.New("target has to be a Func")
	}

	if replacement.Kind() != reflect.Func {
		return 0, nil, errors.New("replacement has to be a Func")
	}

	targetPointer := target.Pointer()
	if _, ok := patches[targetPointer]; ok {
		unpatchValue(targetPointer)
	}

	bytes, originFunc, jumpData, err := replaceFunction(targetPointer, (uintptr)(getPtr(replacement)), replacement.Pointer(), trampoline)
	if err != nil {
		if strings.Contains(err.Error(), "already patched") {
			if p, ok := patches[targetPointer]; ok {
				showInst("origin bytes", targetPointer, p.originalBytes, logger.WarningLevel)
			}
		}
		return 0, nil, err
	}
	patches[targetPointer] = patch{bytes, targetPointer, &replacement, originFunc}
	redirect[originFunc] = targetPointer

	return originFunc, jumpData, nil
}

// PatchPtr 直接将函数跳转的新函数
// 此方式为经过函数签名检查,可能会导致栈帧无法对其导致堆栈调用异常，因此不安全请谨慎使用
// targetPtr 原始函数地址
// replacement 代理函数
func PatchPtr(targetPtr uintptr, replacement interface{}) (*PatchGuard, error) {
	return PatchPtrTrampoline(targetPtr, replacement, nil)
}

// PatchPtrTrampoline 直接将函数跳转的新函数(指定跳板函数)
// 此方式为经过函数签名检查,可能会导致栈帧无法对其导致堆栈调用异常，因此不安全请谨慎使用
// targetPtr 原始函数地址
// replacement 代理函数
// trampoline 跳板函数地址(可不指定,传nil)
func PatchPtrTrampoline(targetPtr uintptr, replacement, trampoline interface{}) (*PatchGuard, error) {
	if p, ok := patches[targetPtr]; ok {
		unpatch(targetPtr, p)
	}

	replacementVal := reflect.ValueOf(replacement)
	trampolinePtr, err := getTrampolinePtr(trampoline)
	if err != nil {
		return nil, err
	}

	bytes, originFunc, jumpData, err := replaceFunction(targetPtr, (uintptr)(getPtr(replacementVal)), replacementVal.Pointer(), trampolinePtr)
	if err != nil {
		return nil, err
	}

	patches[targetPtr] = patch{bytes, targetPtr, &replacementVal, originFunc}
	ptrholder[targetPtr] = replacement
	ptrholder[trampolinePtr] = trampoline
	redirect[originFunc] = targetPtr

	return &PatchGuard{targetPtr, replacementVal, originFunc, jumpData, false}, nil
}

// PatchPtr2Ptr 直接将函数跳转的新函数
// 此方式为经过函数签名检查,可能会导致栈帧无法对其导致堆栈调用异常，因此不安全请谨慎使用
// targetPtr 原始函数地址
// replacement 代理函数跳转地址
// proxy 代理函数地址
// trampoline 跳板函数地址(可不指定,传0)
func PatchPtr2Ptr(targetPtr, replacement, proxy, trampoline uintptr) (*PatchGuard, error) {
	if p, ok := patches[targetPtr]; ok {
		unpatch(targetPtr, p)
	}
	bytes, originFunc, jumpData, err := replaceFunction(targetPtr, replacement, proxy, trampoline)
	if err != nil {
		return nil, err
	}
	patches[targetPtr] = patch{bytes, targetPtr, nil, originFunc}
	redirect[originFunc] = targetPtr
	ptrholder[targetPtr] = replacement

	return &PatchGuard{targetPtr, reflect.ValueOf(nil), originFunc, jumpData, false}, nil
}

// Unpatch removes any monkey patches on target
// returns whether target was patched in the first place
func Unpatch(target interface{}) bool {
	return unpatchValue(reflect.ValueOf(target).Pointer())
}

// PatchInstanceMethod replaces an instance method methodName for the type target with replacement
// Replacement should expect the receiver (of type target) as the first argument
func PatchInstanceMethod(target reflect.Type, methodName string, replacement interface{}) (*PatchGuard, error) {
	return PatchInstanceMethodTrampoline(target, methodName, replacement, nil)
}

// PatchInstanceMethod replaces an instance method methodName for the type target with replacement
// Replacement should expect the receiver (of type target) as the first argument
func PatchInstanceMethodTrampoline(target reflect.Type, methodName string, replacement interface{},
	trampoline interface{}) (*PatchGuard, error) {
	m, ok := target.MethodByName(methodName)
	if !ok {
		return nil, errors.New(fmt.Sprintf("unknown method %s", methodName))
	}
	r := reflect.ValueOf(replacement)
	originFunc, jumpData, err := patchValue(m.Func, r, trampoline)
	if err != nil {
		return nil, err
	}
	ptrholder[m.Func.Pointer()] = replacement

	return &PatchGuard{m.Func.Pointer(), r, originFunc, jumpData, false}, nil
}

// UnpatchAll removes all applied monkeypatches
func UnpatchAll() {
	for target, p := range patches {
		unpatch(target, p)
		delete(patches, target)
		delete(redirect, p.originPtr)
	}
}

// Unpatch removes a monkeypatch from the specified function
// returns whether the function was patched in the first place
func unpatchValue(target uintptr) bool {
	patch, ok := patches[target]
	if !ok {
		return false
	}
	unpatch(target, patch)
	delete(patches, target)
	delete(redirect, patch.originPtr)
	return true
}

func unpatch(target uintptr, p patch) {
	_ = copyToLocation(target, p.originalBytes)
	ShowInst(fmt.Sprintf("unpatch copy to 0x%x", target), target, 20, logger.DebugLevel)
}

// UnpatchInstanceMethod removes the patch on methodName of the target
// returns whether it was patched in the first place
func UnpatchInstanceMethod(target reflect.Type, methodName string) bool {
	m, ok := target.MethodByName(methodName)
	if !ok {
		logger.LogDebugf(fmt.Sprintf("unknown method %s", methodName))
		return false
	}
	return unpatchValue(m.Func.Pointer())
}