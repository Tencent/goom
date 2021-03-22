// Package errobj 收拢了所有的错误类型,错误溯源模板类型等
package errobj

// Cause 带原因的异常类型
type Cause interface {
	// Cause 获取错误的原因
	Cause() error
}

// UnWrapCause 异常转述解包
func UnWrapCause(err error) error {
	if c, ok := err.(Cause); ok {
		return c.Cause()
	}
	return nil
}
