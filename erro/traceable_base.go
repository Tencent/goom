package erro

// TraceableError 可跟踪的错误，异常转述包装
type TraceableError struct {
	err    error
	errStr string
	cause  error
}

// Error 错误描述
func (w *TraceableError) Error() string {
	s := ""
	if w.err != nil {
		s += w.err.Error()
	}
	if w.errStr != "" {
		s += w.errStr
	}
	if w.cause != nil {
		s = s + "\ncause by: " + w.cause.Error()
	}
	return s
}

// Cause 获取错误的原因
func (w *TraceableError) Cause() error {
	return w.cause
}

// NewTraceableErrors 通过 string 描述创建可跟踪的错误
func NewTraceableErrors(errStr string) error {
	return &TraceableError{
		errStr: errStr,
	}
}

// NewTraceableErrorc 通过 string, cause 描述创建可跟踪的错误
func NewTraceableErrorc(errStr string, cause error) error {
	return &TraceableError{
		errStr: errStr,
		cause:  cause,
	}
}

// NewTraceableError 通过 error, cause 描述创建可跟踪的错误
func NewTraceableError(err error, cause error) error {
	return &TraceableError{
		err:   err,
		cause: cause,
	}
}
