// Package iface 生成 interface 的实例对象, 通过 fake iface 结果获得
package iface

import (
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/bytecode"
	"git.code.oa.com/goom/mocker/internal/bytecode/memory"
	"git.code.oa.com/goom/mocker/internal/bytecode/stub"
	"git.code.oa.com/goom/mocker/internal/logger"
)

// interfaceJumpDataLen 默认接口跳转数据长度, 经验数值, 一般采用 icache line 长度(12)的倍数
const interfaceJumpDataLen = 48

// MakeMethodCaller 构造 interface 的方法调用并放到 stub 区
// to 桩函数跳转到的地址
func MakeMethodCaller(to unsafe.Pointer) (uintptr, error) {
	placeholder, _, err := stub.Acquire(interfaceJumpDataLen)
	if err != nil {
		return 0, err
	}

	code := jmpWithRdx(uintptr(to))
	if err := memory.WriteTo(placeholder, code); err != nil {
		return 0, err
	}
	bytecode.PrintInst("gen stub", placeholder, bytecode.PrintMiddle, logger.DebugLevel)
	return placeholder, nil
}

// MakeMethodCallerWithCtx 构造 interface 的方法调用并放到 stub 区
// ctx make Func 对象的上下文地址,即 @see reflect.makeFuncImpl
// to 桩函数最终跳转到另一个地址
func MakeMethodCallerWithCtx(ctx unsafe.Pointer, to uintptr) (uintptr, error) {
	placeholder, _, err := stub.Acquire(interfaceJumpDataLen)
	if err != nil {
		return 0, err
	}

	code := jmpWithRdxAndCtx(uintptr(ctx), to, placeholder)
	if err := memory.WriteTo(placeholder, code); err != nil {
		return 0, err
	}

	bytecode.PrintInst("gen stub", placeholder, bytecode.PrintLong, logger.DebugLevel)
	bytecode.PrintInst("jump to", to, bytecode.PrintLong, logger.DebugLevel)
	return placeholder, nil
}
