// Package proxy_test 对 proxy 包的测试
package proxy_test

import (
	"fmt"
	"runtime/debug"

	"git.woa.com/goom/mocker/internal/logger"
)

//go:noinline
func Caller(i int) int {
	if i <= 0 {
		return 1
	}

	return i * Caller(i-1)
}

//go:noinline
func Caller1(i int) int {
	if i <= 0 {
		return 1
	}

	return i
}

//go:noinline
func Caller2(i int) int {
	i++
	return i
}

// nolint
// Arg 测试参数
type Arg struct {
	field1 string
	field2 map[string]int
	inner  InnerArg
}

//go:noinline
func Caller3(arg Arg) int {
	return len(arg.field2)
}

//go:noinline
func Caller4Eval() {
	arg := &Arg{
		field1: field1,
		field2: map[string]int{},
		inner: InnerArg{
			field1: field1,
			field2: make([]string, 0),
			field3: &InnerField{field3: "ok"},
		},
	}
	Caller4(arg)
}

//go:noinline
func Caller4(arg *Arg) {
	if len(arg.field2) > 0 {
		fmt.Println(len(arg.field2), arg.field1)
	}

	if len(arg.inner.field2) > 0 {
		fmt.Println(len(arg.inner.field2), arg.inner.field1)
	}
}

//go:noinline
func Caller5() int {
	logger.Trace(string(debug.Stack()))
	return 0
}

//go:noinline
func Caller6(a int) func() int {
	return func() int {
		return a + 1
	}
}

//go:noinline
func Caller7(i int) {
	logger.Trace("Caller 7 called")
}

// nolint
// Result 测试参数
type Result struct {
	i     int
	inner *InnerResult
	m     map[string]int
}

// nolint
// InnerResult 测试参数
type InnerResult struct {
	j int
}

//go:noinline
func Caller8(i int) *Result {
	return &Result{
		i: i,
		inner: &InnerResult{
			j: i * 2,
		},
		m: make(map[string]int, 2),
	}
}

//go:noinline
func Caller9(i int) Result {
	return Result{
		i: i,
		inner: &InnerResult{
			j: i * 2,
		},
		m: make(map[string]int, 2),
	}
}

// nolint
// InnerArg 测试参数
type InnerArg struct {
	field1 string
	field2 []string
	field3 *InnerField
}

// nolint
// InnerField 测试参数
type InnerField struct {
	field3 string
}

// nolint
//go:noinline
func ForceStackExpand(i int) int {
	if i <= 0 {
		return 1
	}

	return i * ForceStackExpand(i-1)
}

// nolint
var field1 = "field1"

// nolint
// TestCase 测试用例类型
type TestCase struct {
	funcName string
	// nolint
	funcDef    interface{}
	eval       func()
	trampoline func() interface{}
	proxy      func(interface{}) interface{}
	// nolint
	makeFunc interface{}
	// nolint
	evalMakeFunc func(makeFunc interface{})
}
