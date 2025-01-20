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

func osReadSymbols(reader io.ReaderAt) (*gosym.Table, error) {
	exe, err := elf.NewFile(reader)
	if err != nil {
		return nil, err
	}
	defer exe.Close()

	sect := exe.Section(".text")
	if sect == nil {
		err = fmt.Errorf("Unable to find ELF .text section")
		return nil, err
	}
	textStart := sect.Addr

	sect = exe.Section(".gopclntab")
	if sect == nil {
		err = fmt.Errorf("Unable to find ELF .gopclntab section")
		return nil, err
	}
	lineTableData, err := sect.Data()
	if err != nil {
		return nil, err
	}

	lineTable := gosym.NewLineTable(lineTableData, textStart)
	symTable, err := gosym.NewTable([]byte{}, lineTable)
	if err != nil {
		return nil, err
	}

	// 进一步查找.gosymtab，用于变量表获取
	symbols, err := exe.Symbols()
	if symbols == nil || err != nil {
		// 查找失败, 返回已有的symTable
		if symTable != nil {
			return symTable, nil
		}
		err = fmt.Errorf("Unable to resolve ELF symbols: %v", err)
		return nil, err
	}

	syms := make([]gosym.Sym, 0, len(symbols))
	for i := range symbols {
		syms = append(syms, gosym.Sym{
			Name:  symbols[i].Name,
			Value: symbols[i].Value,
			Type:  symbols[i].Info,
		})
	}
	symTable.Syms = syms
	return symTable, nil
}
