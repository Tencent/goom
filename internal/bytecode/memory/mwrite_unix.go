//go:build !windows
// +build !windows

package memory

import (
	"syscall"
)

// mProtectCrossPage 获取 page 读写权限
func mProtectCrossPage(addr uintptr, length int, prot int) error {
	pageSize := syscall.Getpagesize()
	for p := PageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := RawAccess(p, pageSize)
		if err := syscall.Mprotect(page, prot); err != nil {
			return err
		}
	}
	return nil
}
