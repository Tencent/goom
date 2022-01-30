// Package stub 负责生成和应用桩函数
package stub

import (
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"

	"git.code.oa.com/goom/mocker/internal/patch"
)

// MakeIfaceCaller 构造生成桩函数并放到.text 区
// to 桩函数跳转到的地址
func MakeIfaceCaller(to unsafe.Pointer) (uintptr, error) {
	// acqure space
	placehlder, _, err := acquireSpace(30)
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

	patch.Debug("gen stub", placehlder, 30, logger.DebugLevel)

	return placehlder, nil
}

// MakeIfaceCallerWithCtx 构造生成桩函数并放到.text 区
// ctx make Func 对象的上下文地址,即 @see reflect.makeFuncImpl
// to 桩函数最终跳转到另一个地址
func MakeIfaceCallerWithCtx(ctx unsafe.Pointer, to uintptr) (uintptr, error) {
	// acqure space
	placehlder, _, err := acquireSpace(30)
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

	patch.Debug("genstub", placehlder, 30, logger.DebugLevel)

	return placehlder, nil
}
