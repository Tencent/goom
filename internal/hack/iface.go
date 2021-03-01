package hack

import "unsafe"

// MaxMethod 支持类型的最大方法数量
const (
	// MaxMethod 支持类型的最大方法数量
	MaxMethod = 999
)

// TODO 不同go版本兼容
// Iface 接口结构
type Iface struct {
	Tab  *Itab
	Data unsafe.Pointer
}

// TODO 不同go版本兼容
// 注意: 最多兼容99个方法数量以内的接口
type Itab struct {
	// nolint
	Inter *uintptr
	// nolint
	Type *uintptr
	// nolint
	hash uint32 // copy of Type.hash. Used for type switches.
	_    [4]byte
	Fun  [MaxMethod]uintptr // variable sized. fun[0]==0 means Type does not implement Inter.
}

// Eface 接口结构
type Eface struct {
	// nolint
	rtype unsafe.Pointer
	Data  unsafe.Pointer
}

// UnpackEFace 取出接口
func UnpackEFace(obj interface{}) *Eface {
	return (*Eface)(unsafe.Pointer(&obj))
}
