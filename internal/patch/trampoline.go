package patch

import (
	"errors"
	"fmt"
	"git.code.oa.com/goom/mocker/internal/logger"
)

// fixOriginFuncToTrampoline 将原始函数拷贝到PlayceHolder函数
// 此方式不需要修正pcvalue, 因此相对较安全
// 因trampoline函数需要指定签名,因此只能用于静态代理
// from 原始函数位置
// trampoline 自定义占位函数位置(注意, 自定义占位函数一定要和原函数相同的函数签名,否则栈帧不一致会导致计算调用堆栈时候抛异常)
// jumpInstSize 跳转指令长度, 用于判断需要修复的最小指令长度
// return 跳板函数(即原函数调用入口指针)
func fixOriginFuncToTrampoline(from uintptr, trampoline uintptr, jumpInstSize int) (uintptr, error) {
	// get origin func size
	funcSize, err := GetFuncSize(64, from, true)
	if err != nil {
		logger.LogError("GetFuncSize error", err)
		funcSize = 1024
	}

	logger.LogDebug("origin func size is", funcSize)

	if jumpInstSize >= funcSize {
		ShowInst("origin inst > ", from, 20, logger.InfoLevel)
		return 0, errors.New(fmt.Sprintf("jumpInstSize[%d] is bigger than origin FuncSize[%d], please add your origin func code", jumpInstSize, funcSize))
	}

	// copy origin function
	fixOrigin := rawMemoryRead(from, funcSize)

	showInst("origin inst >>>>> ", from, fixOrigin[:minSize(30, fixOrigin)], logger.DebugLevel)

	// replace replative address to placehlder
	firstFewIns, replaceSize, err := replaceRelativeAddr(from, fixOrigin, trampoline, funcSize, jumpInstSize, true)
	if err != nil {
		return 0, err
	}

	if len(firstFewIns) < len(fixOrigin) {
		// 追加跳转到原函数指令到修复后指令的末尾
		// append jump back to origin func position where next to the broken instructions
		jumpBackData := jmpToOriginFunctionValue(
			trampoline + uintptr(len(firstFewIns)),
			from + (uintptr(replaceSize)))
		fixOrigin = append(firstFewIns, jumpBackData...)
	}

	// get trampoline func size
	trampolineFuncSize, err := GetFuncSize(64, trampoline, false)
	if err != nil {
		logger.LogError("Get trampoline FuncSize error", err)
		return 0, errors.New("Get trampoline FuncSize error:" + err.Error())
	}

	logger.LogDebug("trampoline func size is", trampolineFuncSize)


	if len(fixOrigin) > trampolineFuncSize {
		logger.LogErrorf("fixOriginSize[%d] is bigger than trampoline FuncSize[%d], please add your trampoline func code", len(fixOrigin), trampolineFuncSize)
		ShowInst("trampoline inst > ", trampoline, 35, logger.InfoLevel)
		return 0, errors.New(fmt.Sprintf("fixOriginSize[%d] is bigger than trampoline FuncSize[%d], please add your trampoline func code", len(fixOrigin), trampolineFuncSize))
	} else {
		ShowInst("trampoline inst > ", trampoline, 35, logger.DebugLevel)
	}


	showInst("fixed inst >>>>> ", trampoline, fixOrigin, logger.DebugLevel)

	if err := copyToLocation(trampoline, fixOrigin); err != nil {
		return 0, err
	}
	ShowInst(fmt.Sprintf("tramp copy to 0x%x", trampoline), trampoline, 30, logger.DebugLevel)

	logger.LogDebugf("copy to trampoline %x ", trampoline)

	return trampoline, nil
}
