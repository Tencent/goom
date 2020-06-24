package patch

import (
	"errors"
	"fmt"
	"git.code.oa.com/goom/mocker/internal/logger"
	"reflect"
	"runtime/debug"
	"sync"
	"unsafe"
)


var memoryAccessLock sync.RWMutex

// ReplaceApply 函数调用指针替换执行器
type ReplaceApply func()

// rawMemoryAccess 内存数据读取(非线程安全的)
func rawMemoryAccess(p uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  length,
		Cap:  length,
	}))
}


// rawMemoryRead 内存数据读取(线程安全的)
func rawMemoryRead(p uintptr, length int) []byte {
	memoryAccessLock.RLock()
	defer memoryAccessLock.RUnlock()

	data := rawMemoryAccess(p, length)
	duplucate := make([]byte, length, length)
	copy(duplucate, data)
	return duplucate
}


// from is a pointer to the actual function
// to is a pointer to a go funcvalue
// trampoline 跳板函数地址, 不传递用0表示
func replaceFunction(from, to, proxy, trampoline uintptr) (original []byte, originFunc uintptr, jumpData []byte, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logger.LogErrorf("replaceFunction from=%d to=%d trampoline=%d error:%s", from, to, trampoline, err1)
			logger.LogError(string(debug.Stack()))
			err = err1.(error)
		}
	}()

	logger.LogInfof("starting replace func from=0x%x to=0x%x proxy=0x%x trampoline=0x%x ...", from, to, proxy, trampoline)

	ShowInst("show proxy inst >>>>> ", proxy, 30, logger.DebugLevel)

	// 构造跳转到代理函数的指令
	jumpData = jmpToFunctionValue(from, to)
	// 保存原始指令
	original = rawMemoryRead(from, len(jumpData))
	// 判断是否已经被patch过
	if original[0] == NOP_OPCODE {
		err = errors.New(fmt.Sprintf("from:0x%x is already patched", from))
		return
	}

	// 检测是否支持自动分配跳板函数
	if  trampoline > 0 {
		// 通过跳板函数实现回调原函数
		originFunc, err = fixOriginFuncToTrampoline(from, trampoline, len(jumpData))
		if err != nil {
			return
		}
	}

	return original, originFunc, jumpData, nil
}
