//go:build darwin
// +build darwin

package unexports2

import (
	"debug/gosym"
	"debug/macho"
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

func osReadSymbols(reader io.ReaderAt) (*gosym.Table, error) {
	exe, err := macho.NewFile(reader)
	if err != nil {
		return nil, err
	}
	defer exe.Close()

	var sect *macho.Section
	if sect = exe.Section("__text"); sect == nil {
		err = fmt.Errorf("Unable to find Mach-O __text section")
		return nil, err
	}
	textStart := sect.Addr

	if sect = exe.Section("__gopclntab"); sect == nil {
		err = fmt.Errorf("Unable to find Mach-O __gopclntab section")
		return nil, err
	}
	lineTableData, err := sect.Data()
	if err != nil {
		return nil, err
	}

	lineTable := gosym.NewLineTable(lineTableData, textStart)
	symTable, err := gosym.NewTable([]byte{}, lineTable)
	if err != nil {
		return symTable, err
	}

	if exe.Symtab == nil {
		return symTable, err
	}
	syms := make([]gosym.Sym, 0, len(exe.Symtab.Syms))
	for i := range exe.Symtab.Syms {
		syms = append(syms, gosym.Sym{
			Name:  exe.Symtab.Syms[i].Name,
			Value: exe.Symtab.Syms[i].Value,
			Type:  exe.Symtab.Syms[i].Type,
		})
	}
	symTable.Syms = syms
	return symTable, err
}
