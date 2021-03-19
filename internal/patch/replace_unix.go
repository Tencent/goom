//+build !windows

package patch

import (
	"syscall"
)

var (
	// nolint
	defaultFuncPrologue32 = []byte{0x65, 0x8b, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x8b, 0x89, 0xfc, 0xff, 0xff, 0xff}
	defaultFuncPrologue64 = []byte{0x65, 0x48, 0x8b, 0x0c, 0x25, 0x30, 0x00, 0x00, 0x00, 0x48}
)

// mprotectCrossPage
func mprotectCrossPage(addr uintptr, length int, prot int) {
	pageSize := syscall.Getpagesize()

	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)

		err := syscall.Mprotect(page, prot)
		if err != nil {
			panic("access mem error:" + err.Error())
		}
	}
}

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func CopyToLocation(location uintptr, data []byte) error {
	memoryAccessLock.Lock()
	defer memoryAccessLock.Unlock()

	f := rawMemoryAccess(location, len(data))
	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_EXEC)

	return nil
}
