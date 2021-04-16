package patch

import (
	"fmt"
	"reflect"

	"git.code.oa.com/goom/mocker/internal/logger"
)

// Patch 将函数调用指定代理函数
// target 原始函数
// replacement 代理函数
func Patch(target, replacement interface{}) (*PatchGuard, error) {
	return PatchTrampoline(target, replacement, nil)
}

// PatchTrampoline 将函数调用指定代理函数
// target 原始函数
// replacement 代理函数
// trampoline 指定跳板函数(可不指定,传nil)
func PatchTrampoline(target, replacement interface{}, trampoline interface{}) (*PatchGuard, error) {
	patch := &patch{
		target:      target,
		replacement: replacement,
		trampoline:  trampoline,
	}

	err := patch.patch()
	if err != nil {
		return nil, err
	}

	return patch.Guard(), nil
}

// UnsafePatch 未受类型检查的patch
// target 原始函数
// replacement 代理函数
func UnsafePatch(target, replacement interface{}) (*PatchGuard, error) {
	return UnsafePatchTrampoline(target, replacement, nil)
}

// UnsafePatchTrampoline 未受类型检查的patch
// target 原始函数
// replacement 代理函数
// trampoline 指定跳板函数(可不指定,传nil)
func UnsafePatchTrampoline(target, replacement interface{}, trampoline interface{}) (*PatchGuard, error) {
	patch := &patch{
		target:      target,
		replacement: replacement,
		trampoline:  trampoline,

		targetValue:      reflect.ValueOf(target),
		replacementValue: reflect.ValueOf(replacement),
	}

	err := patch.unsafePatchValue()
	if err != nil {
		return nil, err
	}

	return patch.Guard(), nil
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
	patch := &patch{
		replacement: replacement,
		trampoline:  trampoline,

		replacementValue: reflect.ValueOf(replacement),

		targetPtr: targetPtr,
	}

	err := patch.unsafePatchPtr()
	if err != nil {
		return nil, err
	}

	return patch.Guard(), nil
}

// PatchInstanceMethod replaces an instance method methodName for the type target with replacementValue
// Replacement should expect the receiver (of type target) as the first argument
func PatchInstanceMethod(target reflect.Type, methodName string, replacement interface{}) (*PatchGuard, error) {
	return PatchInstanceMethodTrampoline(target, methodName, replacement, nil)
}

// PatchInstanceMethod replaces an instance method methodName for the type target with replacementValue
// Replacement should expect the receiver (of type target) as the first argument
func PatchInstanceMethodTrampoline(target reflect.Type, methodName string, replacement interface{},
	trampoline interface{}) (*PatchGuard, error) {
	m, ok := target.MethodByName(methodName)
	if !ok {
		return nil, fmt.Errorf("unknown method %s", methodName)
	}

	patch := &patch{
		replacement: replacement,
		trampoline:  trampoline,

		targetValue:      m.Func,
		replacementValue: reflect.ValueOf(replacement),
	}

	err := patch.patchValue()
	if err != nil {
		return nil, err
	}

	return patch.Guard(), nil
}

// Unpatch removes any monkey patches on target
// returns whether target was patched in the first place
func Unpatch(target interface{}) bool {
	return unpatchValue(reflect.ValueOf(target).Pointer())
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

// UnpatchAll removes all applied monkeypatches
func UnpatchAll() {
	for target, p := range patches {
		p.unpatch()
		delete(patches, target)
	}
}

// unpatchValue removes a monkeypatch from the specified function
// returns whether the function was patched in the first place
func unpatchValue(target uintptr) bool {
	p, ok := patches[target]
	if !ok {
		return false
	}

	p.unpatch()
	delete(patches, target)

	return true
}
