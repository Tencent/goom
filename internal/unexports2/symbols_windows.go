//go:build windows
// +build windows

package unexports2

import (
	// "bytes"
	"debug/gosym"
	"debug/pe"
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
	exe, err := pe.NewFile(reader)
	if err != nil {
		return
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
		return
	}

	sect := exe.Section(".text")
	if sect == nil {
		err = fmt.Errorf("Unable to find PE .text section")
		return
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
		err = fmt.Errorf("Could not find PE runtime.pclntab or runtime.epclntab")
		return
	}
	sectionIndex := lineTableStart.SectionNumber - 1
	if sectionIndex < 0 || int(sectionIndex) >= len(exe.Sections) {
		err = fmt.Errorf("Invalid PE format: invalid section number %v", lineTableStart.SectionNumber)
		return
	}
	lineTableData, err := exe.Sections[sectionIndex].Data()
	if err != nil {
		return
	}
	if int(lineTableStart.Value) > len(lineTableData) ||
		int(lineTableEnd.Value) > len(lineTableData) ||
		lineTableStart.Value > lineTableEnd.Value {
		err = fmt.Errorf("Invalid PE pcln start/end indices: %v, %v", lineTableStart.Value, lineTableEnd.Value)
		return
	}
	lineTableData = lineTableData[lineTableStart.Value:lineTableEnd.Value]

	lineTable := gosym.NewLineTable(lineTableData, textStart)
	return gosym.NewTable([]byte{}, lineTable)
}
