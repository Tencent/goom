package hack

import "unsafe"

const (
	// MaxMethod 支持类型的最大方法数量
	MaxMethod = 999
)

// TODO 不同go版本兼容
type Iface struct {
	Tab  *Itab
	Data unsafe.Pointer
}

// TODO 不同go版本兼容
// 注意: 最多兼容99个方法数量以内的接口
type Itab struct {
	// nolint
	inter *uintptr
	// nolint
	_type *uintptr
	// nolint
	hash uint32 // copy of _type.hash. Used for type switches.
	_    [4]byte
	Fun  [MaxMethod]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type Eface struct {
	// nolint
	rtype unsafe.Pointer
	Data  unsafe.Pointer
}

func UnpackEFace(obj interface{}) *Eface {
	return (*Eface)(unsafe.Pointer(&obj))
}
