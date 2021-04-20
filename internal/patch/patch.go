// Package patch 生成指令跳转(到代理函数)并替换.text区内存
// 对于trampoline模式的使用场景，本包实现了指令移动后的修复
package patch

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"git.code.oa.com/goom/mocker/internal/logger"
)

var (
	// lock patches map和内存指令锁
	lock = sync.Mutex{}
	// patches patch缓存
	patches = make(map[uintptr]*patch)
)

// patch is an applied patch
// needed to undo a patch
type patch struct {
	target      interface{}
	replacement interface{}
	trampoline  interface{}

	targetValue      reflect.Value
	replacementValue reflect.Value

	targetPtr      uintptr
	replacementPtr uintptr
	originFuncPtr  uintptr
	trampolinePtr  uintptr

	originalBytes []byte
	jumpBytes     []byte

	guard *Guard
}

// patchValue 对value进行应用代理
func (p *patch) patch() error {
	p.targetValue = reflect.ValueOf(p.target)
	p.replacementValue = reflect.ValueOf(p.replacement)

	return p.patchValue()
}

// patchValue 对value进行应用代理
func (p *patch) patchValue() error {
	// 参数对齐校验 modified by @jake
	checkSignature(p.targetValue.Type(), p.replacementValue.Type())

	return p.unsafePatchValue()
}

// unsafePatchValue 不做类型检查
func (p *patch) unsafePatchValue() error {
	if p.targetValue.Kind() != reflect.Func {
		return errors.New("target has to be a ExportFunc")
	}

	if p.replacementValue.Kind() != reflect.Func {
		return errors.New("replacementValue has to be a ExportFunc")
	}

	targetPointer := p.targetValue.Pointer()
	p.targetPtr = targetPointer

	return p.unsafePatchPtr()
}

// unsafePatchPtr 不做类型检查
func (p *patch) unsafePatchPtr() error {

	replacementPointer := p.replacementValue.Pointer()
	p.replacementPtr = replacementPointer

	if p.trampoline != nil {
		trampolinePtr, err := getTrampolinePtr(p.trampoline)
		if err != nil {
			return err
		}
		p.trampolinePtr = trampolinePtr
	}

	return p.replaceFunc()
}

// replaceFunc 替换函数
func (p *patch) replaceFunc() error {
	// 保证patch和Apply原子性
	Lock()
	defer Unlock()

	if _, ok := patches[p.targetPtr]; ok {
		unpatchValue(p.targetPtr)
	}

	patches[p.targetPtr] = p

	bytes, originFunc, jumpData, err :=
		replaceFunction(p.targetPtr, (uintptr)(getPtr(p.replacementValue)), p.replacementPtr, p.trampolinePtr)
	if err != nil {
		if strings.Contains(err.Error(), "already patched") {
			if pc, ok := patches[p.targetPtr]; ok {
				Debugf("origin bytes", pc.targetPtr, pc.originalBytes, logger.WarningLevel)
			}
		}

		return err
	}

	p.originalBytes = bytes
	p.originFuncPtr = originFunc
	p.jumpBytes = jumpData
	return nil
}

// unpatch do unpatch by uintptr
func (p *patch) unpatch() {
	p.Guard().Unpatch()
}

// restore re patch by target uintptr
// TODO test
func (p *patch) restore() {
	p.Guard().Restore()
}

// Guard 获取PatchGuard
func (p *patch) Guard() *Guard {
	if p.guard != nil {
		return p.guard
	}
	p.guard = &Guard{p.targetPtr,
		p.originFuncPtr,
		p.jumpBytes,
		p.originalBytes,
		false}
	return p.guard
}

// Lock 锁定patches map和内存指令读写
func Lock() {
	lock.Lock()
}

// Unlock 解锁
func Unlock() {
	lock.Unlock()
}
