package hack

import "unsafe"

// TODO 不同go版本兼容
type Iface struct {
	Tab  *Itab
	Data unsafe.Pointer
}

// TODO 不同go版本兼容
type Itab struct {
	inter *uintptr
	_type *uintptr
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	Fun   [3]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type Eface struct {
	rtype unsafe.Pointer
	Data  unsafe.Pointer
}

func UnpackEFace(obj interface{}) *Eface {
	return (*Eface)(unsafe.Pointer(&obj))
}
