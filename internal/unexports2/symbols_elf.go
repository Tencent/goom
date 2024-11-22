//go:build !windows && !darwin
// +build !windows,!darwin

package unexports2

import (
	"debug/elf"
	"debug/gosym"
	"fmt"
	"io"
	"os"
)

func osReadSymbolsFromExeFile() (symTable *gosym.Table, err error) {
	var exePath string
	if exePath, err = os.Executable(); err != nil {
		symTableLoadError = err
		return
	}

	var reader io.ReaderAt
	if reader, err = os.Open(exePath); err != nil {
		symTableLoadError = err
		return
	}

	return osReadSymbols(reader)
}

func osReadSymbols(reader io.ReaderAt) (symTable *gosym.Table, err error) {
	exe, err := elf.NewFile(reader)
	if err != nil {
		return
	}
	defer exe.Close()

	sect := exe.Section(".text")
	if sect == nil {
		err = fmt.Errorf("Unable to find ELF .text section")
		return
	}
	textStart := sect.Addr

	sect = exe.Section(".gopclntab")
	if sect == nil {
		err = fmt.Errorf("Unable to find ELF .gopclntab section")
		return
	}
	lineTableData, err := sect.Data()
	if err != nil {
		return
	}

	lineTable := gosym.NewLineTable(lineTableData, textStart)
	return gosym.NewTable([]byte{}, lineTable)
}
