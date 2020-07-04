package patch

import (
	"syscall"
	"unsafe"
)

// 注意: 此版本暂时不能完整支持windows

const PAGE_EXECUTE_READWRITE = 0x40

var (
	defaultFuncPrologue32 = []byte{0x65, 0x8b, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x8b, 0x89, 0xfc, 0xff, 0xff, 0xff}
	defaultFuncPrologue64 = []byte{0x65, 0x48, 0x8B, 0x0C, 0x25, 0x28, 0x00, 0x00, 0x00}
)

var procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

func virtualProtect(lpAddress uintptr, dwSize int, flNewProtect uint32, lpflOldProtect unsafe.Pointer) error {
	ret, _, _ := procVirtualProtect.Call(
		lpAddress,
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(lpflOldProtect))
	if ret == 0 {
		return syscall.GetLastError()
	}
	return nil
}

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) error {
	f := rawMemoryAccess(location, len(data))

	var oldPerms uint32
	err := virtualProtect(location, len(data), PAGE_EXECUTE_READWRITE, unsafe.Pointer(&oldPerms))
	if err != nil {
		panic(err)
	}
	copy(f, data[:])

	// VirtualProtect requires you to pass in a pointer which it can write the
	// current memory protection permissions to, even if you don't want them.
	var tmp uint32
	return virtualProtect(location, len(data), oldPerms, unsafe.Pointer(&tmp))
}
