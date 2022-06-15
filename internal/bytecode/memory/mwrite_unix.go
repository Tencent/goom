//go:build !windows
// +build !windows

package memory

import (
	"fmt"
	"syscall"

	"git.code.oa.com/goom/mocker/internal/logger"
)

// accessMemGuide access mem error solution guide
const accessMemGuide = "https://iwiki.woa.com/pages/viewpage.action?pageId=1405108952"

// mProtectCrossPage 获取 page 读写权限
func mProtectCrossPage(addr uintptr, length int, prot int) {
	pageSize := syscall.Getpagesize()
	for p := PageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := RawAccess(p, pageSize)
		err := syscall.Mprotect(page, prot)
		if err != nil {
			logger.Consolef(logger.ErrorLevel, "access mem error:permission denied, see details at %s", accessMemGuide)
			panic(fmt.Errorf("access mem error: %w", err))
		}
	}
}
