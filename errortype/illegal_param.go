package errortype

// IllegalParam 参数错误异常
type IllegalParam struct {
	paramName  string
	paramValue string
	funcName   string
}

func (i *IllegalParam) Error() string {
	if len(i.funcName) > 0 {
		return "Illegal param error when call" + i.funcName + ", param=" + i.paramName + ", value=" + i.paramValue
	}

	return "Illegal param error, param=" + i.paramName + ", value=" + i.paramValue
}

// NewIllegalParamError 创建参数异常
// paramName 参数名
// paramValue 参数值
func NewIllegalParamError(paramName string, paramValue string) error {
	return &IllegalParam{paramName: paramName, paramValue: paramValue}
}

// NewIllegalCallError 创建参数异常
// funcName 函数名
// paramName 参数名
// paramValue 参数值
func NewIllegalCallError(funcName string, paramName string, paramValue string) error {
	return &IllegalParam{funcName: funcName, paramName: paramName, paramValue: paramValue}
}
