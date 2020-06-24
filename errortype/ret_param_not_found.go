package errortype

import "strconv"

// ArgNotFound 返回值未找到异常
type ReturnParamNotFound struct {
	funcName string
	arg int
}

func (e *ReturnParamNotFound) Error() string {
	return "arg not found:" + e.funcName + ":" + strconv.Itoa(e.arg)
}

// NewReturnParamNotFoundError 函数未找到
// funcName 函数名称
// index 返回值下标
func NewReturnParamNotFoundError(funcName string, index int) error {
	return &ArgNotFound{funcName:funcName, arg:index}
}
