// +build go1.18

// Package hack 对 go 系统包的 hack, 包含一些系统结构体的 copy，需要和不同的 go 版本保持同步
package hack

import (
	"unsafe"
	_ "unsafe" // 匿名引入
)

// InterceptCallerSkip 拦截器 callerskip
const InterceptCallerSkip = 5

// Firstmoduledata keep async with runtime.Firstmoduledata
//go:linkname Firstmoduledata runtime.firstmoduledata
var Firstmoduledata Moduledata

// Moduledata keep async with runtime.Moduledata
type Moduledata struct {
	pcHeader     *uintptr
	funcnametab  []byte
	cutab        []uint32
	filetab      []byte
	pctab        []byte
	Pclntable    []byte
	Ftab         []Functab
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr
	rodata                uintptr
	gofunc                uintptr // go.func.*

	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []*uintptr

	ptab []interface{}

	pluginpath string
	pkghashes  []interface{}

	modulename   string
	modulehashes []interface{}

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask Bitvector

	typemap map[typeOff]*interface{} // offset to *_rtype in previous module

	bad bool // module failed to load and should be ignored

	Next *Moduledata
}

// Functab Functab
type Functab struct {
	Entry   uint32
	Funcoff uint32
}

// nolint
type textsect struct {
	// nolint
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

// Bitvector Bitvector
type Bitvector struct {
	// nolint
	n int32 // # of bits
	// nolint
	bytedata *uint8
}

// nolint
type typeOff int32 // offset to an *rtype

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
