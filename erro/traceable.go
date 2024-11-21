// Package erro 收拢了所有的错误类型,错误溯源类型等
// 1. 支持打印不同该类型的错误的详细提示, 方便确定问题的解决方法
// 2. 异常可以溯源, 方便问题排查
package erro

// Traceable 带原因的异常类型
// TODO later 可以使用go自带的 %w 功能替代?
type Traceable interface {
	// Cause 获取错误的原因
	Cause() error
}

// Cause 获取错误原因，err对象需要实现Traceable接口才能获取并返回cause error, 否则一律返回false
func Cause(err error) error {
	if c, ok := err.(Traceable); ok {
		return c.Cause()
	}
	return nil
}

// CauseBy 判断错误是否由指定的异常引起,err对象需要实现Traceable接口才能被追踪, 否则一律返回false
// err 判断目标异常
// traceAbleError 指定的异常原因
func CauseBy(err error, traceAbleError Traceable) bool {
	for c := err; c != nil; c = Cause(c) {
		if c == nil {
			return false
		}
		if t, ok := c.(Traceable); ok && t == traceAbleError {
			return true
		}
	}
	return false
}
