package hack

import (
	"unsafe"
	_ "unsafe"
)

// TODO 兼容不同go版本
//go:linkname Firstmoduledata runtime.firstmoduledata
var Firstmoduledata Moduledata

type Moduledata struct {
	Pclntable    []byte
	Ftab         []Functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []textsect
	// Original type was []*_type
	typelinks []int32
	itablinks []*uintptr

	ptab []interface{}

	pluginpath string
	pkghashes  []interface{}

	modulename string
	// Original type was []modulehash
	modulehashes []interface{}

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask Bitvector

	typemap map[typeOff]*interface{} // offset to *_rtype in previous module

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
	n        int32 // # of bits
	bytedata *uint8
}

type textsect struct {
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

type typeOff int32 // offset to an *rtype

// TODO 不同go版本兼容
type Value struct {
	Typ  *uintptr
	Ptr  unsafe.Pointer
	Flag uintptr
}
