// Package memory 包负责内存读写控制
package memory

import (
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

// memoryAccessLock .text 区内存操作度协作
var memoryAccessLock sync.RWMutex

// PageStart page start of memory
func PageStart(addr uintptr) uintptr {
	return addr & ^(uintptr(syscall.Getpagesize() - 1))
}

// RawAccess 内存数据读取(非线程安全的)
// nolint
func RawAccess(addr uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: addr,
		Len:  length,
		Cap:  length,
	}))
}

// RawRead 内存数据读取(线程安全的)
func RawRead(addr uintptr, length int) []byte {
	memoryAccessLock.RLock()
	defer memoryAccessLock.RUnlock()

	data := RawAccess(addr, length)
	duplicate := make([]byte, length)
	copy(duplicate, data)
	return duplicate
}
