package stub

import (
	"fmt"
	"reflect"

	"git.woa.com/goom/mocker/internal/bytecode"
	"git.woa.com/goom/mocker/internal/bytecode/memory"
	"git.woa.com/goom/mocker/internal/logger"
)

const spaceLen = 128

var iCacheHolderAddr uintptr

// ICachePaddingLeft ClearICache 左侧占位
func ICachePaddingLeft()

// ClearICache 汇编函数声明: 清理 icache 缓存
func ClearICache()

func init() {
	iCacheHolderAddr = reflect.ValueOf(ClearICache).Pointer()
	// 兼容 go 1.17(1.17以上会对 assembler 函数进行 wrap, 需要找到其内部的调用)
	innerAddr, err := bytecode.GetInnerFunc(64, iCacheHolderAddr)
	if innerAddr > 0 && err == nil {
		iCacheHolderAddr = innerAddr
	}
	offset := reflect.ValueOf(ICachePaddingLeft).Pointer()
	logger.Debugf("icache func init success: %x", offset)
}

// WriteICacheFn 写入 icache clear 函数数据
func WriteICacheFn(data []byte) (uintptr, error) {
	s, err := acquireICacheFn()
	if err != nil {
		return 0, err
	}
	switch s.typ {
	case TypeMMap:
		copy(*s.Space, data[:])
		return s.Addr, nil
	case TypeHolder:
		return s.Addr, memory.WriteToNoFlushNoLock(s.Addr, data)
	default:
		return 0, fmt.Errorf("ICacheFn write fail, illegal type: %d", s.typ)
	}
}

// acquireICacheFn 获取 icache 执行空间
func acquireICacheFn() (*Space, error) {
	if addr, space, err := acquireFromMMap(spaceLen); err == nil {
		return &Space{
			Addr:  addr,
			Space: space,
			typ:   TypeMMap,
		}, nil
	}

	return &Space{
		Addr:  iCacheHolderAddr,
		Space: nil,
		typ:   TypeHolder,
	}, nil
}
