package errortype

// IllegalStatus 状态错误异常
type IllegalStatus struct {
	funcName string
	msg      string
}

// IllegalStatus 状态错误异常
func (i *IllegalStatus) Error() string {
	return "Illegal status error when call " + i.funcName + " msg:" + i.msg
}

// NewIllegalStatusError 状态参数异常
// funcName 函数名
// msg 信息
func NewIllegalStatusError(funcName string, msg string) error {
	return &IllegalStatus{funcName: funcName, msg: msg}
}
