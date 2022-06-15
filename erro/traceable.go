// Package erro 收拢了所有的错误类型,错误溯源类型等
// 1. 支持打印不同该类型的错误的详细提示, 方便确定问题的解决方法
// 2. 异常可以溯源, 方便问题排查
package erro

// Traceable 带原因的异常类型
// TODO later 可以使用go自带的 %w 功能替代
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
