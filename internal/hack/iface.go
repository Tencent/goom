// Package hack 对 go 系统包的 hack, 包含一些系统结构体的 copy，需要和每个 go 版本保持同步
package hack

import "unsafe"

const (
	// MaxMethod 支持类型的最大方法数量
	MaxMethod = 999
)

// Iface 接口结构
// TODO 不同 go 版本兼容
type Iface struct {
	// Tab 为接口类型的方法表
	Tab *Itab
	// Data 为接口变量所持有的对实现类型接收体的地址
	Data unsafe.Pointer
}

// Itab keeps sync with runtime.itab
// TODO 不同 go 版本兼容
// 注意: 最多兼容99个方法数量以内的接口
type Itab struct {
	// nolint
	Inter *uintptr
	// nolint
	Type *uintptr
	// nolint
	hash uint32 // copy of Type.hash. Used for type switches.
	_    [4]byte
	// Fun 为方法表映射、排序同接口方法定义的顺序
	Fun [MaxMethod]uintptr // variable sized. fun[0]==0 means Type does not implement Inter.
}

// Eface 接口结构
type Eface struct {
	// nolint
	rtype unsafe.Pointer
	// Data 为 interface{}类型变量指向的 Iface 类型变量的地址
	Data unsafe.Pointer
}

// UnpackEFace 取出接口对象
func UnpackEFace(obj interface{}) *Eface {
	return (*Eface)(unsafe.Pointer(&obj))
}
