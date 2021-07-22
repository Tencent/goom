package erro

import "strings"

const prefix = "func not found: "

// FuncNotFound 函数未找到异常
type FuncNotFound struct {
	funcName    string
	suggestions []string
}

// FuncNotFound 函数未找到异常
func (e *FuncNotFound) Error() string {
	msg := prefix + e.funcName
	if e.suggestions == nil {
		return msg
	}

	var noEmptyStrings []string
	for _, v := range e.suggestions {
		if v == "" {
			continue
		}
		noEmptyStrings = append(noEmptyStrings, v)
	}
	if len(noEmptyStrings) == 0 {
		return msg
	}

	msg = prefix + "\n  " + e.funcName
	msg += ", do you mean? \n* "
	msg += strings.Join(noEmptyStrings, "\n* ")
	msg += "\n （〃＾∇＾）oo"
	return msg
}

// NewFuncNotFoundError 函数未找到
// funcName 函数名称
func NewFuncNotFoundError(funcName string) error {
	return &FuncNotFound{funcName: funcName}
}

// NewFuncNotFoundErrorWithSuggestion 函数未找到并给出提示
// funcName 函数名称
func NewFuncNotFoundErrorWithSuggestion(funcName string, suggestions []string) error {
	return &FuncNotFound{funcName: funcName, suggestions: suggestions}
}
