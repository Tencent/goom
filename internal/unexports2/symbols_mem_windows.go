package unexports2

import (
	"log"
	"reflect"
	"syscall"
)

func osGetProcessBaseAddress() (address uintptr) {
	osInitProcess()
	// This will be 0x400000 unless ASLR is on.

	// https://docs.microsoft.com/en-us/previous-versions/bb985992(v=msdn.10)?redirectedfrom=MSDN
	address, _, _ = getModuleHandle.Call(0)
	return
}

func osGetPageSize() int {
	// https://docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info
	var memory [20]uint32
	ptr := reflect.ValueOf(&memory).Elem().UnsafeAddr()
	getSystemInfo.Call(ptr)
	pageSize := int(memory[1]) // dwPageSize
	if pageSize <= 0 {
		log.Printf("go-subvert: Warning: GetSystemInfo returned page size of %v\n", pageSize)
		pageSize = 0x1000
	}
	return pageSize
}

var kernel32 *syscall.LazyDLL
var virtualProtect *syscall.LazyProc
var getModuleHandle *syscall.LazyProc
var getSystemInfo *syscall.LazyProc

func osInitProcess() {
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	virtualProtect = kernel32.NewProc("VirtualProtect")
	virtualProtect.Addr() // Forces a panic if not found
	getModuleHandle = kernel32.NewProc("GetModuleHandleA")
	getModuleHandle.Addr()
	getSystemInfo = kernel32.NewProc("GetSystemInfo")
	getSystemInfo.Addr()
}
