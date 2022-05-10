package patch

import (
	"fmt"
	"reflect"
	deb "runtime/debug"
	"sync"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"
)

// memoryAccessLock .text 区内存操作度协作
var memoryAccessLock sync.RWMutex

// nolint
// rawMemoryAccess 内存数据读取(非线程安全的)
func rawMemoryAccess(ptr uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: ptr,
		Len:  length,
		Cap:  length,
	}))
}

// rawMemoryRead 内存数据读取(线程安全的)
func rawMemoryRead(ptr uintptr, length int) []byte {
	memoryAccessLock.RLock()
	defer memoryAccessLock.RUnlock()

	data := rawMemoryAccess(ptr, length)
	duplicate := make([]byte, length)
	copy(duplicate, data)

	return duplicate
}

// replaceFunction 在函数 from 里面, 织入对 to 的调用指令，同时将 from 织入前的指令恢复至 trampoline 这个地址
// from is a pointer to the actual function
// to is a pointer to a go function value
// trampoline 跳板函数地址, 不传递用0表示
func replaceFunction(from, to, proxy, trampoline uintptr) (original []byte, originFunc uintptr,
	jumpData []byte, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logger.LogErrorf("replaceFunction from=%d to=%d trampoline=%d error:%s", from, to, trampoline, err1)
			logger.LogError(string(deb.Stack()))

			var ok bool

			err, ok = err1.(error)
			if !ok {
				err = fmt.Errorf("%s", err1)
			}
		}
	}()

	logger.LogInfof("starting replace func from=0x%x to=0x%x proxy=0x%x trampoline=0x%x ...", from, to, proxy, trampoline)

	Debug("show proxy inst >>>>> ", proxy, 30, logger.DebugLevel)

	// 构造跳转到代理函数的指令
	jumpData = jmpToFunctionValue(from, to)

	// get origin func size
	funcSize, err := GetFuncSize(defaultArchMod, from, false)
	if err != nil {
		logger.LogError("GetFuncSize error", err)

		funcSize = defaultFuncSize
	}

	// 如果需要织入的跳转指令的长度大于原函数指令长度,则任务是无法织入指令
	if len(jumpData) >= funcSize {
		Debug("origin inst > ", from, insSizePrintShort, logger.InfoLevel)
		return nil, 0, nil, fmt.Errorf(
			"jumpInstSize[%d] is bigger than origin FuncSize[%d], cannot do pathes", len(jumpData), funcSize)
	}

	// 保存原始指令
	original = rawMemoryRead(from, len(jumpData))
	// 判断是否已经被 patch 过
	if checkAlreadyPatch(original) {
		err = fmt.Errorf("from:0x%x is already patched", from)
		return
	}

	Debugf("origin >>>>> ", from, rawMemoryRead(from, 30), logger.DebugLevel)

	// 检测是否支持自动分配跳板函数
	if trampoline > 0 {
		// 通过跳板函数实现回调原函数
		originFunc, err = fixOriginFuncToTrampoline(from, trampoline, len(jumpData))
		if err != nil {
			return
		}
	}

	return original, originFunc, jumpData, nil
}
