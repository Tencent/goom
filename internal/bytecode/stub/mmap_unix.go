//go:build !windows
// +build !windows

package stub

import (
	"syscall"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"
)

// acquireFromMMap enough executable space from holder
func acquireFromMMap(len int) (uintptr, *[]byte, error) {
	executableSpace, err := syscall.Mmap(
		-1,
		0,
		len,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_SHARED|syscall.MAP_ANON)
	if err != nil {
		logger.Debugf("acquireFromMMap fail: %v\n", err)
		return 0, nil, err
	}
	addr := (uintptr)(unsafe.Pointer(&executableSpace[0]))
	return addr, &executableSpace, nil
}
