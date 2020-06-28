package errortype

import "strconv"

// ArgNotFound 参数未找到异常
type ArgNotFound struct {
	funcName string
	arg      int
}

func (e *ArgNotFound) Error() string {
	return "arg not found:" + e.funcName + ":" + strconv.Itoa(e.arg)
}

// NewFuncNotFoundError 函数未找到
// funcName 函数名称
// index 参数下标
func NewArgNotFoundError(funcName string, index int) error {
	return &ArgNotFound{funcName: funcName, arg: index}
}
