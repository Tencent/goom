//go:build !windows
// +build !windows

package memory

import (
	"fmt"
	"syscall"

	"github.com/tencent/goom/internal/logger"
)

// accessMemGuide access mem error solution guide
const accessMemGuide = "https://github.com/tencent/goom"

// WriteTo this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func WriteTo(addr uintptr, data []byte) error {
	memoryAccessLock.Lock()
	defer memoryAccessLock.Unlock()

	f := RawAccess(addr, len(data))
	if err := mProtectCrossPage(addr, len(data), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC); err != nil {
		// mac 环境下使用 hack 方式绕过权限检查
		if e := writeTo(addr, data); e == nil {
			return nil
		}
		errorDetail(err)
	}
	copy(f, data[:])
	if err := mProtectCrossPage(addr, len(data), syscall.PROT_READ|syscall.PROT_EXEC); err != nil {
		errorDetail(err)
	}
	return nil
}

func errorDetail(err error) {
	logger.Consolef(logger.ErrorLevel, "access mem error:permission denied, see details at %s", accessMemGuide)
	panic(fmt.Errorf("access mem error: %w", err))
}
