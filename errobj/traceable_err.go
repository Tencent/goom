package errobj

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
		s = s + "\ncause: " + w.cause.Error()
	}
	return s
}

// Traceable 获取错误的原因
func (w *TraceableError) Cause() error {
	return w.cause
}

// NewTraceableError 创建可跟踪的错误
func NewTraceableError(err error, cause error) error {
	return &TraceableError{
		err:   err,
		cause: cause,
	}
}

// NewTraceableErrorf 通过string描述创建可跟踪的错误
func NewTraceableErrorf(errStr string, cause error) error {
	return &TraceableError{
		errStr: errStr,
		cause:  cause,
	}
}
