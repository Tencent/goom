package stub

import (
	"fmt"
	"reflect"

	"git.code.oa.com/goom/mocker/internal/bytecode/memory"
	"git.code.oa.com/goom/mocker/internal/logger"
)

const spaceLen = 128

var iCacheHolderAddr uintptr

// ICachePaddingLeft ClearICache 左侧占位
func ICachePaddingLeft()

// ClearICache 汇编函数声明: 清理 icache 缓存
func ClearICache()

func init() {
	iCacheHolderAddr = reflect.ValueOf(ClearICache).Pointer()
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
	} else {
		return &Space{
			Addr:  iCacheHolderAddr,
			Space: nil,
			typ:   TypeHolder,
		}, nil
	}
}
