//go:build go1.21
// +build go1.21

package hack

import (
	"runtime"
	"unsafe"
	_ "unsafe" // 匿名引入
)

// InterceptCallerSkip 拦截器 callerskip
const InterceptCallerSkip = 5

// Firstmoduledata keep async with runtime.Firstmoduledata
//
//go:linkname Firstmoduledata runtime.firstmoduledata
var Firstmoduledata Moduledata

// nolint
// Moduledata keep async with runtime.Moduledata
type Moduledata struct {
	NotInHeap // Only in static data

	pcHeader     *uintptr
	Funcnametab  []byte
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
	covctrs, ecovctrs     uintptr
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

	// This slice records the initializing tasks that need to be
	// done to start up the program. It is built by the linker.
	inittasks []interface{}

	modulename   string
	modulehashes []interface{}

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask Bitvector

	_ map[typeOff]*interface{} // offset to *_rtype in previous module

	_ bool // module failed to load and should be ignored

	Next *Moduledata
}

type NotInHeap struct{ _ nih }
type nih struct{}

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

// CheckNameOffOverflow check nameOff overflow
func CheckNameOffOverflow(f *runtime.Func, md *Moduledata) bool {
	fc := (*FuncInfo)(unsafe.Pointer(f))
	if fc.NameOff >= int32(len(md.Funcnametab)) {
		return true
	}
	return false
}

// FuncInfo keep async with runtime2.go/type _func struct{}
type FuncInfo struct {
	EntryOff uint32 // start pc, as offset from moduledata.text/pcHeader.textStart
	NameOff  int32  // function name, as index into moduledata.funcnametab.
	// .....
}
