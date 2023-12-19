// Package iface 生成 interface 的实例对象, 通过 fake iface 结果获得
package iface

import (
	"unsafe"

	"github.com/tencent/goom/internal/bytecode"
	"github.com/tencent/goom/internal/bytecode/stub"
	"github.com/tencent/goom/internal/logger"
)

// interfaceJumpDataLen 默认接口跳转数据长度, 经验数值, 一般采用 icache line 长度(12)的倍数
const interfaceJumpDataLen = 48

// MakeMethodCaller 构造 interface 的方法调用并放到 stub 区
// to 桩函数跳转到的地址
func MakeMethodCaller(to unsafe.Pointer) (uintptr, error) {
	space, err := stub.Acquire(interfaceJumpDataLen)
	if err != nil {
		return 0, err
	}

	code := jmpWithRdx(uintptr(to))
	if err := stub.Write(space, code); err != nil {
		return 0, err
	}
	bytecode.PrintInst("gen stub", space.Addr, bytecode.PrintMiddle, logger.DebugLevel)
	return space.Addr, nil
}

// MakeMethodCallerWithCtx 构造 interface 的方法调用并放到 stub 区
// ctx make Func 对象的上下文地址,即 @see reflect.makeFuncImpl
// to 桩函数最终跳转到另一个地址
func MakeMethodCallerWithCtx(ctx unsafe.Pointer, to uintptr) (uintptr, error) {
	space, err := stub.Acquire(interfaceJumpDataLen)
	if err != nil {
		return 0, err
	}

	code := jmpWithRdx(uintptr(ctx))
	if err := stub.Write(space, code); err != nil {
		return 0, err
	}

	bytecode.PrintInst("gen stub", space.Addr, bytecode.PrintLong, logger.DebugLevel)
	bytecode.PrintInst("jump to", to, bytecode.PrintLong, logger.DebugLevel)
	return space.Addr, nil
}
