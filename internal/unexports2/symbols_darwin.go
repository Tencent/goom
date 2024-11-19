//go:build darwin
// +build darwin

package unexports2

import (
	"bytes"
	"debug/gosym"
	"debug/macho"
	"fmt"
	"io"
	"os"
)

func osReadSymbolsFromMemory() (symTable *gosym.Table, err error) {
	if processBaseAddress == 0 {
		return nil, fmt.Errorf("Base address not found")
	}
	reader := bytes.NewReader(SliceAtAddress(processBaseAddress, 0x10000000))
	return osReadSymbols(reader)
}

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
	exe, err := macho.NewFile(reader)
	if err != nil {
		return
	}
	defer exe.Close()

	var sect *macho.Section
	if sect = exe.Section("__text"); sect == nil {
		err = fmt.Errorf("Unable to find Mach-O __text section")
		return
	}
	textStart := sect.Addr

	if sect = exe.Section("__gopclntab"); sect == nil {
		err = fmt.Errorf("Unable to find Mach-O __gopclntab section")
		return
	}
	lineTableData, err := sect.Data()
	if err != nil {
		return
	}

	lineTable := gosym.NewLineTable(lineTableData, textStart)
	return gosym.NewTable([]byte{}, lineTable)
}
