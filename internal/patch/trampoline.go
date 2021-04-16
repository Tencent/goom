package patch

import (
	"errors"
	"fmt"

	"git.code.oa.com/goom/mocker/internal/logger"
)

const (
	// 默认需要修复的函数长度
	defaultFuncSize = 1024
	// 默认系统位数、暂时不支持32位的
	defaultArchMod = 64
	// 日志相关
	// 默认打印的指令数量(短)
	insSizePrintShort = 20
	// 默认打印的指令数量(中)
	insSizePrintMiddle = 30
	// 默认打印的指令数量(长)
	insSizePrintLong = 35
)

// fixOriginFuncToTrampoline 将原始函数from的指令到trampoline指向的地址(在PlaceHolder区内存区段内)
// 此方式不需要修正pcvalue, 因此相对较安全
// 因trampoline函数需要指定签名,因此只能用于静态代理
// from 原始函数位置
// trampoline 自定义占位函数位置(注意, 自定义占位函数一定要和原函数相同的函数签名,否则栈帧不一致会导致计算调用堆栈时候抛异常)
// jumpInstSize 跳转指令长度, 用于判断需要修复的最小指令长度
// return 跳板函数(即原函数调用入口指针)
func fixOriginFuncToTrampoline(from uintptr, trampoline uintptr, jumpInstSize int) (uintptr, error) {
	// get origin func size
	originFuncSize, err := GetFuncSize(defaultArchMod, from, false)
	if err != nil {
		logger.LogError("GetFuncSize error", err)

		originFuncSize = defaultFuncSize
	}

	// get trampoline func size
	trampFuncSize, err := GetFuncSize(defaultArchMod, trampoline, false)
	if err != nil {
		logger.LogError("GetFuncSize error", err)

		trampFuncSize = 20
	}

	logger.LogDebug("origin func size is", originFuncSize)

	// 如果需要修复的指令长度大于trampoline函数指令长度,则任务是无法修复
	if jumpInstSize >= trampFuncSize {
		Debug("origin inst > ", from, insSizePrintShort, logger.InfoLevel)
		return 0, fmt.Errorf(
			"jumpInstSize[%d] is bigger than trampoline FuncSize[%d], "+
				"please fill your trampoline func code", jumpInstSize, originFuncSize)
	}

	// copy origin function
	fixOrigin := rawMemoryRead(from, originFuncSize)

	debug("origin inst >>>>> ", from, fixOrigin[:minSize(insSizePrintMiddle, fixOrigin)], logger.DebugLevel)

	// replace relative address to placeholder
	firstFewIns, replaceSize, err := replaceRelativeAddr(from, fixOrigin, trampoline, originFuncSize, jumpInstSize, true)
	if err != nil {
		return 0, err
	}

	if len(firstFewIns) < len(fixOrigin) {
		// 追加跳转到原函数指令到修复后指令的末尾
		// append jump back to origin func position where next to the broken instructions
		jumpBackData := jmpToOriginFunctionValue(
			trampoline+uintptr(len(firstFewIns)),
			from+(uintptr(replaceSize)))
		fixOrigin = append(firstFewIns, jumpBackData...)
	}

	// get trampoline func size
	trampolineFuncSize, err := GetFuncSize(defaultArchMod, trampoline, false)
	if err != nil {
		logger.LogError("Get trampoline FuncSize error", err)
		return 0, errors.New("Get trampoline FuncSize error:" + err.Error())
	}

	logger.LogDebug("trampoline func size is", trampolineFuncSize)

	if len(fixOrigin) > trampolineFuncSize {
		logger.LogErrorf("fixOriginSize[%d] is bigger than trampoline FuncSize[%d], please add your "+
			"trampoline func code", len(fixOrigin), trampolineFuncSize)
		Debug("trampoline inst > ", trampoline, insSizePrintLong, logger.InfoLevel)

		return 0, fmt.Errorf("fixOriginSize[%d] is bigger than trampoline FuncSize[%d], "+
			"please add your trampoline func code", len(fixOrigin), trampolineFuncSize)
	}

	Debug("trampoline inst > ", trampoline, insSizePrintLong, logger.DebugLevel)
	debug("fixed inst >>>>> ", trampoline, fixOrigin, logger.DebugLevel)

	if err := CopyToLocation(trampoline, fixOrigin); err != nil {
		return 0, err
	}

	Debug(fmt.Sprintf("tramp copy to 0x%x", trampoline), trampoline, insSizePrintMiddle, logger.DebugLevel)
	logger.LogDebugf("copy to trampoline %x ", trampoline)

	return trampoline, nil
}
