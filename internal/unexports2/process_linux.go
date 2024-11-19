package unexports2

import (
	"bufio"
	"os"
	"strconv"
	"syscall"
)

func osGetProcessBaseAddress() (address uintptr) {
	path := "/proc/self/maps"

	fileReader, err := os.Open(path)
	if err != nil {
		return
	}
	reader := bufio.NewReader(fileReader)
	addrString, err := reader.ReadBytes('-')
	if err != nil {
		return
	}

	addrString = addrString[:len(addrString)-1]
	parsed, err := strconv.ParseUint(string(addrString), 16, 64)
	if err != nil {
		return
	}
	address = uintptr(parsed)
	return
}

func osGetPageSize() int {
	return syscall.Getpagesize()
}

func osInitProcess() {
	// Nothing to do
}
