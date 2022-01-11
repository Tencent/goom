//go:build !windows
// +build !windows

package patch

import (
	"syscall"

	"git.code.oa.com/goom/mocker/internal/logger"
)

var (
	//nolint build by build tag
	// defaultFuncPrologue32 32位系统function Prologue
	defaultFuncPrologue32 = []byte{0x65, 0x8b, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x8b, 0x89, 0xfc, 0xff, 0xff, 0xff}
	// defaultFuncPrologue64 64位系统function Prologue
	defaultFuncPrologue64 = []byte{0x65, 0x48, 0x8b, 0x0c, 0x25, 0x30, 0x00, 0x00, 0x00, 0x48}
	// arm64 func prologue
	armFuncPrologue64 = []byte{0x81, 0x0B, 0x40, 0xF9, 0xE2, 0x83, 0x00, 0xD1, 0x5F, 0x00, 0x01, 0xEB}
	// accessMemGuide access mem error solution guide
	accessMemGuide = "https://iwiki.woa.com/pages/viewpage.action?pageId=1405108952"
)

// mprotectCrossPage 获取page读写权限
func mprotectCrossPage(addr uintptr, length int, prot int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Log2Consolef(logger.ErrorLevel, "access mem error:permission denied, see details at %s", accessMemGuide)
			panic(err)
		}
	}()
	pageSize := syscall.Getpagesize()

	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)

		err := syscall.Mprotect(page, prot)
		if err != nil {
			panic("access mem error:" + err.Error())
		}
	}
}

// CopyToLocation this function is super unsafe
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
