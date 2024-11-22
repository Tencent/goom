package unexports2

import (
	"debug/gosym"
	"fmt"
)

var (
	symTable          *gosym.Table
	symTableLoadError error
)

func loadSymbolTable() (table *gosym.Table, err error) {
	if err = symTableLoadError; err != nil {
		return
	}

	if table = symTable; table != nil {
		return
	}

	table, err = osReadSymbolsFromExeFile()
	symTableLoadError = err
	if err != nil {
		symTable = nil
	} else if table == nil {
		err = fmt.Errorf("Unknown error: symbol table was nil")
	} else {
		symTable = table
	}

	symTable = table
	return
}

// GetFunctionSymbol returns the symbols for a given function.
func GetFunctionSymbol(function interface{}) (symbol *gosym.Func, err error) {
	var table *gosym.Table
	if table, err = GetSymbolTable(); err != nil {
		return
	}

	address, err := getFunctionAddress(function)
	if err != nil {
		return
	}
	symbol = table.PCToFunc(uint64(address))
	if symbol == nil {
		err = fmt.Errorf("Function symbol at %x not found", address)
	}
	return
}

func getFunctionSymbolByName(name string) (symbol *gosym.Func, err error) {
	var table *gosym.Table
	if table, err = GetSymbolTable(); err != nil {
		return
	}

	symbol = table.LookupFunc(name)
	if symbol == nil {
		err = fmt.Errorf("%v: function symbol not found", name)
	}
	return
}
