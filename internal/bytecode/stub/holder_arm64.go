package stub

import (
	"reflect"
)

var addr uintptr

// ClearICache 汇编函数声明: 清理 icache 缓存
func ClearICache()

func init() {
	addr = reflect.ValueOf(ClearICache).Pointer()
}

// ICacheHolder 获取 icache 地址
func ICacheHolder() uintptr {
	return addr
}
