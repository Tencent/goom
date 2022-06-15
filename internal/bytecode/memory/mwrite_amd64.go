//go:build !windows
// +build !windows

package memory

import (
	"syscall"
)

// WriteTo this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func WriteTo(addr uintptr, data []byte) error {
	memoryAccessLock.Lock()
	defer memoryAccessLock.Unlock()

	f := RawAccess(addr, len(data))
	mProtectCrossPage(addr, len(data), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
	copy(f, data[:])
	mProtectCrossPage(addr, len(data), syscall.PROT_READ|syscall.PROT_EXEC)
	return nil
}
