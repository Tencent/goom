//go:build go1.21
// +build go1.21

package hack

import (
	"unsafe"
	_ "unsafe" // 匿名引入
)

// InterceptCallerSkip 拦截器 callerskip
const InterceptCallerSkip = 5

type NotInHeap struct{ _ nih }
type nih struct{}

// Functab Functab
type Functab struct {
	Entry   uint32
	Funcoff uint32
}

// Bitvector Bitvector
type Bitvector struct {
	// nolint
	n int32 // # of bits
	// nolint
	bytedata *uint8
}

// Func convenience struct for modifying the underlying code pointer of a function
// value. The actual struct has other values, but always starts with a code
// pointer.
type Func struct {
	CodePtr uintptr
}

// Value reflect.Value
type Value struct {
	Typ  *uintptr
	Ptr  unsafe.Pointer
	Flag uintptr
}

// FuncInfo keep async with runtime2.go/type _func struct{}
type FuncInfo struct {
	EntryOff uint32 // start pc, as offset from moduledata.text/pcHeader.textStart
	NameOff  int32  // function name, as index into moduledata.funcnametab.
	// .....
}
