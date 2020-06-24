package errortype

// FuncNotFound 函数未找到异常
type FuncNotFound struct {
	funcName string
}

func (e *FuncNotFound) Error() string {
	return "func not found:" + e.funcName
}

// NewFuncNotFoundError 函数未找到
// funcName 函数名称
func NewFuncNotFoundError(funcName string) error {
	return &FuncNotFound{funcName:funcName}
}

