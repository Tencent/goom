// Package arg 负责参数表达式构造和执行, 执行结果用于 When 参数匹配(Matcher)
package arg

// Pair 参数和返回值键值对
type Pair struct {
	// Params 参数列表
	Params interface{}
	// Return 返回值
	Return interface{}
}
