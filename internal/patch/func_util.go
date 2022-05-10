package patch

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/unexports"
	"git.code.oa.com/goom/mocker/internal/x86asm"
)

var (
	// funcSizeCache 函数长度缓存
	funcSizeCache = make(map[uintptr]int)
	// funcSizeReadLock 并发读写 funcSizeCache 锁
	funcSizeReadLock sync.Mutex
)

// pageStart page start of memory
func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

// GetInnerFunc Get the first real func location from wrapper
// not absolutely safe
func GetInnerFunc(mode int, start uintptr) (uintptr, error) {
	prologueLen := len(funcPrologue)
	code := rawMemoryRead(start, 16) // instruction takes at most 16 bytes

	int3Found := false
	curLen := 0

	for {
		inst, err := x86asm.Decode(code, mode)
		if err != nil || (inst.Opcode == 0 && inst.Len == 1 && inst.Prefix[0] == x86asm.Prefix(code[0])) {
			return 0, nil
		}

		if inst.Len == 1 && code[0] == 0xcc {
			int3Found = true
		} else if int3Found {
			return 0, nil
		}

		if inst.Op.String() == callInsName {
			relativeAddr := decodeRelativeAddr(&inst, code, inst.PCRelOff)
			return start + (uintptr)(relativeAddr) + uintptr(inst.Len), nil
		}

		curLen = curLen + inst.Len
		code = rawMemoryRead(start+uintptr(curLen), 16) // instruction takes at most 16 bytes

		if bytes.Equal(funcPrologue, code[:prologueLen]) {
			return 0, nil
		}
	}
}

// value value
type value struct {
	_   uintptr
	ptr unsafe.Pointer
}

// getPtr 获取函数的调用地址(和函数的指令地址不一样)
func getPtr(v reflect.Value) unsafe.Pointer {
	return (*value)(unsafe.Pointer(&v)).ptr
}

// isNil 判断 interface{}是否为空
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).Elem().IsNil()
	}

	return reflect.ValueOf(i).IsNil()
}

// getTrampolinePtr 获取跳板函数的地址
func getTrampolinePtr(trampoline interface{}) (uintptr, error) {
	trampolinePtr := uintptr(0)

	if !isNil(trampoline) {
		trampolineType := reflect.TypeOf(trampoline)

		if trampolineType.Kind() == reflect.Ptr {
			trampolinePtr = reflect.ValueOf(trampoline).Elem().Pointer()
		} else if trampolineType.Kind() == reflect.Func {
			trampolinePtr = reflect.ValueOf(trampoline).Pointer()
		}

		logger.LogDebugf("trampoline value: 0x%x 0x%x", getPtr(reflect.ValueOf(trampoline)), trampolinePtr)
	}

	return trampolinePtr, nil
}

// IsPtr 判断 interface{}是否为指针类型
func IsPtr(value interface{}) bool {
	if value == nil {
		return false
	}

	t := reflect.TypeOf(value)

	return t.Kind() == reflect.Ptr
}

// Debug Debug 调试内存指令替换,对原指令、替换之后的指令进行输出对比
func Debug(name string, from uintptr, size int, level int) {
	_, funcName, _ := unexports.FindFuncByPtr(from)
	instBytes := rawMemoryRead(from, size)
	Debugf(fmt.Sprintf("show [%s = %s] inst>>: ", name, funcName), from, instBytes, level)
}

// Debugf 调试内存指令替换,对原指令、替换之后的指令进行输出对比
func Debugf(title string, from uintptr, copyOrigin []byte, level int) {
	if logger.LogLevel < level {
		return
	}

	logger.LogImportant(title)

	startAddr := (uint64)(from)

	for pos := 0; pos < len(copyOrigin); {
		// read 16 bytes at most each time
		endPos := pos + 16
		if endPos > len(copyOrigin) {
			endPos = len(copyOrigin)
		}

		code := copyOrigin[pos:endPos]
		ins, err := x86asm.Decode(code, 64)

		if err != nil {
			logger.LogImportantf("[0] 0x%x: inst decode error:%s", startAddr+(uint64)(pos), err)

			if ins.Len == 0 {
				pos = pos + 1
			} else {
				pos = pos + ins.Len
			}

			continue
		}

		if ins.Opcode == 0 {
			if ins.Len == 0 {
				pos = pos + 1
			} else {
				pos = pos + ins.Len
			}

			continue
		}

		if ins.PCRelOff <= 0 {
			logger.LogImportantf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s", ins.Len,
				startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(code[:ins.Len]))

			pos = pos + ins.Len

			continue
		}

		isAdd := true

		for i := 0; i < len(ins.Args); i++ {
			arg := ins.Args[i]
			if arg == nil {
				break
			}

			addrArgs := arg.String()
			if strings.HasPrefix(addrArgs, ".-") || strings.Contains(addrArgs, "RIP-") {
				isAdd = false
			}
		}

		offset := pos + ins.PCRelOff

		relativeAddr := decodeAddress(copyOrigin[offset:offset+ins.PCRel], ins.PCRel)
		if !isAdd && relativeAddr > 0 {
			relativeAddr = -relativeAddr
		}

		logger.LogImportantf("[%d] 0x%x:\t%s\t\t%-30s\t\t%s\t\tabs:0x%x", ins.Len,
			startAddr+(uint64)(pos), ins.Op, ins.String(), hex.EncodeToString(code[:ins.Len]),
			from+uintptr(pos)+uintptr(relativeAddr)+uintptr(ins.Len))

		pos = pos + ins.Len
	}
}

// minSize 最小 size，不超出 fixOrigin 长度的 size 大小
func minSize(showSize int, fixOrigin []byte) int {
	if showSize > len(fixOrigin) {
		showSize = len(fixOrigin)
	}

	return showSize
}

// checkSignature 检测两个函数类型的参数的内存区段是否一致
func checkSignature(targetType reflect.Type, replacementType reflect.Type) bool {
	// 检测参数对齐
	if targetType.NumIn() != replacementType.NumIn() {
		panic(fmt.Sprintf("func signature mismatch, args len must:%d, actual:%d",
			targetType.NumIn(), replacementType.NumIn()))
	}
	if targetType.NumOut() != replacementType.NumOut() {
		panic(fmt.Sprintf("func signature mismatch, returns len must:%d, actual:%d",
			targetType.NumOut(), replacementType.NumOut()))
	}
	for i := 0; i < targetType.NumIn(); i++ {
		if targetType.In(i).Size() != replacementType.In(i).Size() {
			panic(fmt.Sprintf("func signature mismatch, args %d's size must:%d, actual:%d",
				i, targetType.In(i).Size(), replacementType.In(i).Size()))
		}
	}
	for i := 0; i < targetType.NumOut(); i++ {
		if targetType.Out(i).Size() != replacementType.Out(i).Size() {
			panic(fmt.Sprintf("func signature mismatch, returns %d's size must:%d, actual:%d",
				i, targetType.Out(i).Size(), replacementType.Out(i).Size()))
		}
	}
	return true
}
