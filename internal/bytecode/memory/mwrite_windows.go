package memory

import (
	"fmt"
	"syscall"
	"unsafe"
)

// 注意: 此版本暂时不能完整支持 windows
// pageExecuteReadwrite page 窗口大小
const pageExecuteReadwrite = 0x40

var procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

// virtualProtect 获取 page 读写权限
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

// WriteTo this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func WriteTo(addr uintptr, data []byte) error {
	memoryAccessLock.Lock()
	defer memoryAccessLock.Unlock()

	f := RawAccess(addr, len(data))

	var oldPerms uint32
	err := virtualProtect(addr, len(data), pageExecuteReadwrite, unsafe.Pointer(&oldPerms))
	if err != nil {
		panic(fmt.Errorf("access mem error: %w", err))
	}
	copy(f, data[:])

	// VirtualProtect requires you to pass in a pointer which it can write the
	// current memory protection permissions to, even if you don't want them.
	var tmp uint32
	return virtualProtect(addr, len(data), oldPerms, unsafe.Pointer(&tmp))
}
