package proxy

import (
	"fmt"
	"runtime/debug"

	"git.code.oa.com/goom/mocker/internal/logger"
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
	i += 1
	return i
}

type Arg struct {
	field1 string
	field2 map[string]int
	inner  InnerArg
}

var argPtr = &Arg{
	field1: field1,
	field2: nil,
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
	//fmt.Println("called4", arg.field1)
	return
}

//go:noinline
func Caller5() int {
	logger.LogTrace(string(debug.Stack()))
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
	logger.LogTrace("Caller 7 called")
}

type Result struct {
	i     int
	inner *InnerResult
	m     map[string]int
}

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

type InnerArg struct {
	field1 string
	field2 []string
	field3 *InnerField
}

type InnerField struct {
	field3 string
}

//go:noinline
func ForceStackExpand(i int) int {
	if i <= 0 {
		return 1
	}
	return i * ForceStackExpand(i-1)
}

var field1 = "field1"

type TestCase struct {
	funcName     string
	funcDef      interface{}
	eval         func()
	trampoline   func() interface{}
	proxy        func(interface{}) interface{}
	makefunc     interface{}
	evalMakeFunc func(makefunc interface{})
}
