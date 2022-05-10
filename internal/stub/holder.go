// Package stub 负责生成和应用桩函数
package stub

import (
	"errors"
	"reflect"
	"sync/atomic"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/patch"

	"git.code.oa.com/goom/mocker/internal/logger"
)

// placeHolderIns 占位实例
var placeHolderIns *PlaceHolder

// PlaceHolder 占位对象
type PlaceHolder struct {
	// count hook 次数统计
	count int
	// off 当前占位函数使用的偏移量
	off uintptr
	// min 占位函数起始位置
	min uintptr
	// max 占位函数末尾位置
	max uintptr
}

// Placeholder 汇编函数声明: 占位函数
func Placeholder()

func init() {
	offset := reflect.ValueOf(Placeholder).Pointer()

	// 兼容 go 1.17(1.17以上会对 assembler 函数进行 wrap, 需要找到其内部的调用)
	innerOffset, err := patch.GetInnerFunc(64, offset)
	if innerOffset > 0 && err == nil {
		offset = innerOffset
	}

	size, err := patch.GetFuncSize(64, offset, false)
	if err != nil {
		logger.LogError("GetFuncSize error", err)

		size = 102400
	}

	placeHolderIns = &PlaceHolder{
		count: 0,
		off:   offset,
		min:   offset,
		max:   uintptr(size) + offset,
	}

	logger.LogDebugf("Placeholder pointer: %d %d\n", placeHolderIns.min, offset)
}

// addOff add mapping offset to origin func
func addOff(from uintptr, used uintptr) error {
	// add up to off
	newOffset := atomic.AddUintptr(&placeHolderIns.off, used+16)
	if newOffset+used > placeHolderIns.max {
		logger.LogError("placeholder space usage overflow", placeHolderIns.count, "hook functions")
		return errors.New("placeholder space usage overflow")
	}

	logger.LogDebug("add offset map, size:", used)

	return nil
}

// acquireSpace check if has enough holder space
func acquireSpace(funcSize int) (uintptr, []byte, error) {
	placeholder := atomic.LoadUintptr(&placeHolderIns.off)
	if placeholder+uintptr(funcSize) > placeHolderIns.max {
		logger.LogError("placeholder space usage overflow")
		return 0, nil, errors.New("placeholder space usage overflow")
	}

	bytes := (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: placeholder,
		Len:  funcSize,
		Cap:  funcSize,
	}))

	return placeholder, *bytes, nil
}
