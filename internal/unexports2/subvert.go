// Package subvert provides functions to subvert go's type & memory protections,
// and expose unexported values & functions.
//
// This is not a power to be taken lightly! It's expected that you're fully
// versed in how the go type system works, and why there are protections and
// restrictions in the first place. Using this package incorrectly will quickly
// lead to undefined behavior and bizarre crashes, even segfaults or nuclear
// missile launches.

// YOU HAVE BEEN WARNED!
package unexports2

import (
	"debug/gosym"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

// MakeWritable clears a value's RO flags. The RO flags are generally used to
// determine whether a value is exported (and thus accessible) or not.
func MakeWritable(v *reflect.Value) error {
	if !rvFlagsFound {
		return rvFlagsError
	}
	*getRVFlagPtr(v) &= ^rvFlagRO
	return nil
}

// MakeAddressable adds the addressable flag to a value, allowing you to take
// its address. The most common reason for making an object non-addressable is
// because it's allocated on the stack or in read-only memory.
//
// Making a pointer to a stack value will cause undefined behavior if you
// attempt to access it outside of the stack-allocated object's scope.
//
// Do not write to an object in read-only memory. It would be bad.
func MakeAddressable(v *reflect.Value) error {
	if !rvFlagsFound {
		return rvFlagsError
	}
	*getRVFlagPtr(v) |= rvFlagAddr
	return nil
}

// SliceAtAddress turns a memory range into a go slice.
//
// No checks are made as to whether the memory is writable or even readable.
//
// Do not append to the slice.
func SliceAtAddress(address uintptr, length int) []byte {
	pageSize := uintptr(syscall.Getpagesize())
	dataSize := uintptr(length)
	var (
		errnoResult syscall.Errno
		success     bool
	)
	for p := pageStart(address); p < address+dataSize; p += pageSize {
		_, _, errno := syscall.Syscall(syscall.SYS_MPROTECT, p, pageSize, syscall.PROT_READ|syscall.PROT_EXEC)
		if errno != 0 {
			errnoResult = errno
		}
		if errno == 0 {
			success = true
		}
	}
	if !success {
		panic(fmt.Errorf("access mem error: %w", errnoResult))
	}
	return rawAccess(address, length)
}

func rawAccess(addr uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: addr,
		Len:  length,
		Cap:  length,
	}))
}

// pageStart page start of memory
func pageStart(addr uintptr) uintptr {
	return addr & ^(uintptr(syscall.Getpagesize() - 1))
}

// GetSliceAddr gets the address of a slice
func GetSliceAddr(slice []byte) uintptr {
	pSlice := (*unsafe.Pointer)((unsafe.Pointer)(&slice))
	return uintptr(*pSlice)
}

// ExposeFunction exposes a function or method, allowing you to bypass export
// restrictions. It looks for the symbol specified by funcSymName and returns a
// function with its implementation, or nil if the symbol wasn't found.
//
// funcSymName must be the exact symbol name from the binary. Use AllFunctions()
// to find it. If your program doesn't have any references to a function, it
// will be omitted from the binary during compilation. You can prevent this by
// saving a reference to it somewhere, or calling a function that indirectly
// references it.
//
// templateFunc MUST have the correct function type, or else undefined behavior
// will result!
//
// Example:
//
//	exposed := ExposeFunction("reflect.methodName", (func() string)(nil))
//	if exposed != nil {
//	    f := exposed.(func() string)
//	    fmt.Printf("Result of reflect.methodName: %v\n", f())
//	}
func ExposeFunction(funcSymName string, templateFunc interface{}) (function interface{}, err error) {
	fn, err := getFunctionSymbolByName(funcSymName)
	if err != nil {
		return
	}
	return newFunctionWithImplementation(templateFunc, uintptr(fn.Entry))
}

// GetSymbolTable loads (if necessary) and returns the symbol table for this process
func GetSymbolTable() (*gosym.Table, error) {
	if symTable == nil && symTableLoadError == nil {
		symTable, symTableLoadError = loadSymbolTable()
	}

	return symTable, symTableLoadError
}

// AllFunctions returns the name of every function that has been compiled
// into the current binary. Use it as a debug helper to see if a function
// has been compiled in or not.
func AllFunctions() (functions map[string]bool, err error) {
	var table *gosym.Table
	if table, err = GetSymbolTable(); err != nil {
		return
	}

	functions = make(map[string]bool)
	for _, function := range table.Funcs {
		functions[function.Name] = true
	}
	return
}

func init() {
	initReflectValue()
	initProcess()
}

const is64BitUintptr = uint64(^uintptr(0)) == ^uint64(0)
