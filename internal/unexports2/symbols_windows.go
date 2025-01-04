//go:build windows
// +build windows

package unexports2

import (
	"debug/gosym"
	"debug/pe"
	"fmt"
	"io"
	"os"

	"github.com/tencent/goom/erro"
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
	exe, err := pe.NewFile(reader)
	if err != nil {
		return nil, err
	}
	defer exe.Close()

	var imageBase uint64
	switch oh := exe.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		imageBase = uint64(oh.ImageBase)
	case *pe.OptionalHeader64:
		imageBase = oh.ImageBase
	default:
		err = fmt.Errorf("Unrecognized PE format")
		return nil, err
	}

	sect := exe.Section(".text")
	if sect == nil {
		err = erro.NewTraceableErrorc("Unable to find PE .text section", erro.LdFlags)
		return nil, err
	}
	textStart := imageBase + uint64(sect.VirtualAddress)

	findSymbol := func(symbols []*pe.Symbol, name string) *pe.Symbol {
		for _, s := range symbols {
			if s.Name == name {
				return s
			}
		}
		return nil
	}

	lineTableStart := findSymbol(exe.Symbols, "runtime.pclntab")
	lineTableEnd := findSymbol(exe.Symbols, "runtime.epclntab")
	if lineTableStart == nil || lineTableEnd == nil {
		err = erro.NewTraceableErrorc("Could not find PE runtime.pclntab or runtime.epclntab", erro.LdFlags)
		return nil, err
	}
	sectionIndex := lineTableStart.SectionNumber - 1
	if sectionIndex < 0 || int(sectionIndex) >= len(exe.Sections) {
		err = fmt.Errorf("Invalid PE format: invalid section number %v", lineTableStart.SectionNumber)
		return nil, err
	}
	lineTableData, err := exe.Sections[sectionIndex].Data()
	if err != nil {
		return nil, err
	}
	if int(lineTableStart.Value) > len(lineTableData) ||
		int(lineTableEnd.Value) > len(lineTableData) ||
		lineTableStart.Value > lineTableEnd.Value {
		err = fmt.Errorf("Invalid PE pcln start/end indices: %v, %v", lineTableStart.Value, lineTableEnd.Value)
		return nil, err
	}
	lineTableData = lineTableData[lineTableStart.Value:lineTableEnd.Value]

	lineTable := gosym.NewLineTable(lineTableData, textStart)
	symTable, err := gosym.NewTable([]byte{}, lineTable)
	if err != nil {
		return symTable, err
	}
	if exe.Symbols == nil {
		return symTable, err
	}

	syms := make([]gosym.Sym, 0, len(exe.Symbols))
	for i := range exe.Symbols {
		syms = append(syms, gosym.Sym{
			Name:  exe.Symbols[i].Name,
			Value: uint64(exe.Symbols[i].Value),
			Type:  byte(exe.Symbols[i].Type), // is that correct?
		})
	}
	symTable.Syms = syms
	return symTable, err
}
