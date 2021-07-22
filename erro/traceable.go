// Package erro 收拢了所有的错误类型,错误溯源类型等
package erro

// Traceable 带原因的异常类型
type Traceable interface {
	// Cause 获取错误的原因
	Cause() error
}

// CauseOf 获取错误原因
func CauseOf(err error) error {
	if c, ok := err.(Traceable); ok {
		return c.Cause()
	}
	return nil
}
