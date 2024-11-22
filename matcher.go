package mocker

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"git.woa.com/goom/mocker/arg"
)

// BaseMatcher 参数匹配基类
type BaseMatcher struct {
	results [][]reflect.Value
	curNum  int32
	funTyp  reflect.Type
	// resultsPtr 持有参数指针, 防止被回收
	resultsPtr []interface{}
}

// newBaseMatcher 创建新参数匹配基类
func newBaseMatcher(results []interface{}, funTyp reflect.Type) *BaseMatcher {
	resultVs := make([][]reflect.Value, 0)
	if results != nil {
		// TODO results check
		result, err := arg.I2V(results, outTypes(funTyp))
		if err != nil {
			panic("Return Value (" + fmt.Sprintf("%v", results) + ") error: " + err.Error())
		}
		resultVs = append(resultVs, result)
	}
	return &BaseMatcher{
		results:    resultVs,
		curNum:     0,
		funTyp:     funTyp,
		resultsPtr: results,
	}
}

// Result 回参
func (c *BaseMatcher) Result() []reflect.Value {
	if len(c.results) <= 1 {
		return c.results[c.curNum]
	}

	curNum := atomic.LoadInt32(&c.curNum)
	if length := len(c.results); curNum >= int32(length) {
		return c.results[length-1]
	}

	atomic.AddInt32(&c.curNum, 1)
	return c.results[curNum]
}

// AddResult 添加结果
func (c *BaseMatcher) AddResult(results []interface{}) {
	// TODO results check
	result, err := arg.I2V(results, outTypes(c.funTyp))
	if err != nil {
		panic("Return Value (" + fmt.Sprintf("%v", results) + ") error: " + err.Error())
	}
	c.results = append(c.results, result)
}

// EmptyMatch 没有返回参数的匹配器
type EmptyMatch struct {
	*AlwaysMatcher
}

// newEmptyMatch 创建无参数匹配器
func newEmptyMatch() *EmptyMatch {
	return &EmptyMatch{}
}

// Result 返回参数
func (c *EmptyMatch) Result() []reflect.Value {
	return []reflect.Value{}
}

// DefaultMatcher 参数匹配
// 入参个数必须和函数或方法参数个数一致,
// 比如: When(
//
//	In(3, 4), // 第一个参数是 In
//	Any()) // 第二个参数是 Any
type DefaultMatcher struct {
	*BaseMatcher
	isMethod bool
	exprs    []arg.Expr
}

// newDefaultMatch 创建新参数匹配
func newDefaultMatch(args []interface{}, results []interface{}, isMethod bool, funTyp reflect.Type) *DefaultMatcher {
	e, err := arg.ToExpr(args, inTypes(isMethod, funTyp))
	if err != nil {
		panic(fmt.Sprintf("Call When("+fmt.Sprintf("%v", args)+") error: %v", err))
	}
	return &DefaultMatcher{
		exprs:       e,
		BaseMatcher: newBaseMatcher(results, funTyp),
		isMethod:    isMethod,
	}
}

// Match 判断是否匹配
func (c *DefaultMatcher) Match(args []reflect.Value) bool {
	if c.isMethod {
		args = args[1:]
	}
	if len(args) != len(c.exprs) {
		return false
	}

	for i, expr := range c.exprs {
		v, err := expr.Eval([]reflect.Value{args[i]})
		if err != nil {
			// TODO add mocker and method name to message
			panic(fmt.Sprintf("param[%d] match fail: %v", i, err))
		}
		if !v {
			return false
		}
	}
	return true
}

// ContainsMatcher 包含类型的参数匹配
// 当参数为多个时, In 的每个条件各使用一个数组表示:
// .In([]interface{}{3, Any()}, []interface{}{4, Any()})
type ContainsMatcher struct {
	*BaseMatcher
	expr     *arg.InExpr
	isMethod bool
}

// newContainsMatch 创建新的包含类型的参数匹配
func newContainsMatch(args []interface{}, results []interface{}, isMethod bool,
	funTyp reflect.Type) *ContainsMatcher {
	in := arg.In(args...)
	err := in.Resolve(inTypes(isMethod, funTyp))
	if err != nil {
		// TODO add mocker and method name to message
		panic(fmt.Sprintf("create param match fail: %v", err))
	}
	return &ContainsMatcher{
		expr:        in,
		BaseMatcher: newBaseMatcher(results, funTyp),
		isMethod:    isMethod,
	}
}

// Match 判断是否匹配
func (c *ContainsMatcher) Match(args []reflect.Value) bool {
	if c.isMethod {
		args = args[1:]
	}
	v, err := c.expr.Eval(args)
	if err != nil {
		// TODO add mocker and method name to message
		panic(fmt.Sprintf("param match fail: %v", err))
	}
	return v
}

// AlwaysMatcher 默认匹配
type AlwaysMatcher struct {
	*BaseMatcher
}

// newAlwaysMatch 创建新的默认匹配
func newAlwaysMatch(results []interface{}, funTyp reflect.Type) *AlwaysMatcher {
	if results == nil {
		return nil
	}
	return &AlwaysMatcher{
		BaseMatcher: newBaseMatcher(results, funTyp),
	}
}

// Match 总是匹配
func (c *AlwaysMatcher) Match(_ []reflect.Value) bool {
	return true
}
