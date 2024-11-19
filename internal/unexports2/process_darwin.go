package unexports2

import (
	"bytes"
	"fmt"
	"syscall"
)

func osGetProcessBaseAddress() (address uintptr) {
	defer func() {
		if e := recover(); e != nil {
			address = 0
		}
	}()

	// Do a guess assuming this function's address is close to the base.
	addr, err := getFunctionAddress(osGetProcessBaseAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	startAddress := addr & ^uintptr(0xffffff)
	length := 0x400000

	// Mach-O images start with either FEEDFACE (32-bit) or FEEDFACF (64-bit)
	toFind := []byte{0xcf, 0xfa, 0xed, 0xfe}
	if !is64BitUintptr {
		toFind[0] = 0xce
	}
	asBytes := SliceAtAddress(startAddress, length)
	index := bytes.Index(asBytes, toFind)
	if index < 0 {
		return
	}
	address = startAddress + uintptr(index)
	return
}

func osGetPageSize() int {
	return syscall.Getpagesize()
}

func osInitProcess() {
	// Nothing to do
}
