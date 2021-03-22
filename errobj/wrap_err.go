package errobj

// WrapError 异常转述包装
type WrapError struct {
	err    error
	errStr string
	cause  error
}

func (w *WrapError) Error() string {
	s := ""
	if w.err != nil {
		s += w.err.Error()
	}
	if w.errStr != "" {
		s += w.errStr
	}
	if w.cause != nil {
		s = s + "\ncause: " + w.cause.Error()
	}
	return s
}

// Cause 获取错误的原因
func (w *WrapError) Cause() error {
	return w.cause
}

// NewWrapError 创建异常转述包装对象
func NewWrapError(err error, cause error) error {
	return &WrapError{
		err:   err,
		cause: cause,
	}
}

// NewWrapErrorS 创建异常转述包装对象
func NewWrapErrorS(errStr string, cause error) error {
	return &WrapError{
		errStr: errStr,
		cause:  cause,
	}
}

// UnWrap 异常转述解包
func UnWrap(err error) error {
	if w, ok := err.(*WrapError); ok {
		return w.cause
	}
	return nil
}
