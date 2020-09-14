package hack

import (
	"unsafe"
	// nolint
	_ "unsafe"
)

// Firstmoduledata TODO 兼容不同go版本
//go:linkname Firstmoduledata runtime.firstmoduledata
var Firstmoduledata Moduledata

type Moduledata struct {
	Pclntable []byte
	Ftab      []Functab
	// nolint
	filetab []uint32
	// nolint
	findfunctab uintptr
	// nolint
	minpc, maxpc uintptr
	// nolint
	text, etext uintptr
	// nolint
	noptrdata, enoptrdata uintptr
	// nolint
	data, edata uintptr
	// nolint
	bss, ebss uintptr
	// nolint
	noptrbss, enoptrbss uintptr
	// nolint
	end, gcdata, gcbss uintptr
	// nolint
	types, etypes uintptr
	// nolint
	textsectmap []textsect
	// Original type was []*_type
	// nolint
	typelinks []int32
	// nolint
	itablinks []*uintptr
	// nolint
	ptab []interface{}
	// nolint
	pluginpath string
	// nolint
	pkghashes []interface{}
	// nolint
	modulename string
	// nolint
	// Original type was []modulehash
	modulehashes []interface{}
	// nolint
	hasmain uint8 // 1 if module contains the main function, 0 otherwise
	// nolint
	gcdatamask, gcbssmask Bitvector
	// nolint
	typemap map[typeOff]*interface{} // offset to *_rtype in previous module
	// nolint
	bad bool // module failed to load and should be ignored

	Next *Moduledata
}

type Functab struct {
	Entry   uintptr
	Funcoff uintptr
}

// Convenience struct for modifying the underlying code pointer of a function
// value. The actual struct has other values, but always starts with a code
// pointer.
// TODO 不同go版本兼容
type Func struct {
	CodePtr uintptr
}

type Bitvector struct {
	// nolint
	n int32 // # of bits
	// nolint
	bytedata *uint8
}

// nolint
type textsect struct {
	// nolint
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

// nolint
type typeOff int32 // offset to an *rtype

// TODO 不同go版本兼容
type Value struct {
	Typ  *uintptr
	Ptr  unsafe.Pointer
	Flag uintptr
}