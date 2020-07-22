package stub

import (
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"

	"git.code.oa.com/goom/mocker/internal/patch"
)

// GenStub 构造桩函数
// ctx 上下文地址 @see reflect.makeFuncImpl
// to 桩函数最终跳转到另一个地址
func GenStub(ctx unsafe.Pointer, to uintptr) (uintptr, error) {

	code := jmpWithRdx(uintptr(ctx), to)

	// acqure space
	placehlder, _, err := acqureSpace(len(code))
	if err != nil {
		return 0, err
	}

	if err := patch.CopyToLocation(placehlder, code); err != nil {
		return 0, err
	}

	if err := addOff(0, uintptr(len(code))); err != nil {
		return 0, err
	}

	patch.ShowInst("genstub", placehlder, 30, logger.DebugLevel)

	return placehlder, nil
}

// jmpWithRdx Assembles a jump to a function value
// dx DX寄存器
// to 跳转目标地址
func jmpWithRdx(dx, to uintptr) (value []byte) {
	return []byte{
		0x48, 0x8B,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx,to
		0xFF, 0xD0,     // jmp QWORD PTR [rdx]
	}
}
