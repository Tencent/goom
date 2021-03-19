// Package stub负责生成和应用桩函数
package stub

import (
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"

	"git.code.oa.com/goom/mocker/internal/patch"
)

// MakeIfaceCaller 构造生成桩函数并放到.text区
// to 桩函数跳转到的地址
func MakeIfaceCaller(to unsafe.Pointer) (uintptr, error) {
	// acqure space
	placehlder, _, err := acqureSpace(30)
	if err != nil {
		return 0, err
	}

	code := jmpWithRdx(uintptr(to))

	if err := patch.CopyToLocation(placehlder, code); err != nil {
		return 0, err
	}

	if err := addOff(0, uintptr(len(code))); err != nil {
		return 0, err
	}

	patch.ShowInst("genstub", placehlder, 30, logger.DebugLevel)

	return placehlder, nil
}

// MakeIfaceCallerWithCtx 构造生成桩函数并放到.text区
// ctx make Func对象的上下文地址,即 @see reflect.makeFuncImpl
// to 桩函数最终跳转到另一个地址
func MakeIfaceCallerWithCtx(ctx unsafe.Pointer, to uintptr) (uintptr, error) {
	// acqure space
	placehlder, _, err := acqureSpace(30)
	if err != nil {
		return 0, err
	}

	code := jmpWithRdxAndCtx(uintptr(ctx), to, placehlder)

	if err := patch.CopyToLocation(placehlder, code); err != nil {
		return 0, err
	}

	if err := addOff(0, uintptr(len(code))); err != nil {
		return 0, err
	}

	patch.ShowInst("genstub", placehlder, 30, logger.DebugLevel)

	return placehlder, nil
}

// jmpWithRdx Assembles a jump to a clourse function value
// dx DX寄存器
func jmpWithRdx(dx uintptr) (value []byte) {
	return []byte{
		0x48, 0xBA,
		byte(dx),
		byte(dx >> 8),
		byte(dx >> 16),
		byte(dx >> 24),
		byte(dx >> 32),
		byte(dx >> 40),
		byte(dx >> 48),
		byte(dx >> 56), // movabs rdx,dx
		0xFF, 0x22,     // jmp QWORD PTR [rdx]
	}
}

// jmpWithRdxAndCtx Assembles a jump to a function value
// ctx context地址
// to 跳转目标地址
// from 跳转来源地址
func jmpWithRdxAndCtx(ctx, to, from uintptr) (value []byte) {
	var dis uint32
	if to > from {
		dis = uint32(int32(to-from) - 5)
		dis = dis + 10
	} else {
		dis = uint32(-int32(from-to) - 5)
		dis = dis - 10
	}

	return []byte{
		0x48, 0xBA,
		byte(ctx),
		byte(ctx >> 8),
		byte(ctx >> 16),
		byte(ctx >> 24),
		byte(ctx >> 32),
		byte(ctx >> 40),
		byte(ctx >> 48),
		byte(ctx >> 56), // movabs rdx,ctx

		0xe9,
		byte(dis),
		byte(dis >> 8),
		byte(dis >> 16),
		byte(dis >> 24), // jmp dis
	}
}
