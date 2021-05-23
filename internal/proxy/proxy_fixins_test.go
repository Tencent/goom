package proxy

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"testing"

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
	for j := 0; j < 10; j++ {
		i += j
	}
	return i
}

type Arg struct {
	field1 string
	field2 map[string]int
}

//go:noinline
func Caller3(arg Arg) int {
	//if len(arg.field2) > 0 {
	//	fmt.Println(len(arg.field2))
	//}
	return 2 + len(arg.field1) + len(arg.field2)
}

//go:noinline
func Caller4(arg *Arg) int {
	return 0
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
}

//go:noinline
func Caller8(i int) int {
tag:
	return i
tag1:
	return -i
	if i < 0 {
		goto tag
	} else {
		goto tag1
	}
}

var testCases1 = []struct {
	funcName   string
	funcDef    interface{}
	eval       func(t *testing.T)
	trampoline func() interface{}
	proxy      func(interface{}) interface{}
}{
	{
		funcName: "Caller",
		funcDef:  Caller,
		eval: func(t *testing.T) {
			if r := Caller(5); r != 120 {
				t.Fatalf("want result: %d, real: %d", 120, r)
			}
		},
		trampoline: func() interface{} {
			var result = func(i int) int {
				fmt.Println("trampoline")
				return i + 10
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			return func(i int) int {
				logger.LogTrace("proxy Caller called, args", i)
				originFunc, _ := origin.(*func(i int) int)
				return (*originFunc)(i)
			}
		},
	},
	{
		funcName: "Caller1",
		funcDef:  Caller1,
		eval: func(t *testing.T) {
			if r := Caller1(-1); r != 1 {
				t.Fatalf("want result: %d, real: %d", 1, r)
			}
		},
		trampoline: func() interface{} {
			var result = func(i int) int {
				fmt.Println("trampoline")
				return i + 20
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(i int) int {
				logger.LogTrace("proxy Caller1 called, args", i)
				originFunc, _ := origin1.(*func(i int) int)
				return (*originFunc)(i)
			}
		},
	},
	{
		funcName: "Caller2",
		funcDef:  Caller2,
		eval: func(t *testing.T) {
			if r := Caller2(5); r != 50 {
				t.Fatalf("want result: %d, real: %d", 50, r)
			}

		},
		trampoline: func() interface{} {
			var result = func(i int) int {
				fmt.Println("trampoline")
				return i + 30
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(i int) int {
				logger.LogTrace("proxy Caller2 called, args", i)
				originFunc, _ := origin1.(*func(i int) int)
				return (*originFunc)(i)
			}
		},
	},
	{
		funcName: "Caller3",
		funcDef:  Caller3,
		eval: func(t *testing.T) {
			if r := Caller3(Arg{
				field1: "field1",
				field2: nil,
			}); r != 8 {
				t.Fatalf("want result: %d, real: %d", 8, r)
			}
		},
		trampoline: func() interface{} {
			var result = func(arg Arg) int {
				fmt.Println("trampoline")
				return 40
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(arg Arg) int {
				fmt.Println("Caller3")
				logger.LogTrace("proxy Caller3 called, args", arg)
				originFunc, _ := origin1.(*func(arg Arg) int)
				fmt.Println("Caller3-1")
				result := (*originFunc)(arg)
				fmt.Println("Caller3-2")
				return result
			}
		},
	},
	{
		funcName: "Caller4",
		funcDef:  Caller4,
		eval: func(t *testing.T) {
			if r := Caller4(&Arg{
				field1: "field1",
				field2: make(map[string]int),
			}); r != 0 {
				t.Fatalf("want result: %d, real: %d", 0, r)
			}
		},
		trampoline: func() interface{} {
			var result = func(arg *Arg) int {
				fmt.Println("trampoline")
				return 50
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(arg *Arg) int {
				logger.LogTrace("proxy Caller4 called, args", arg)
				originFunc, _ := origin1.(*func(arg *Arg) int)
				return (*originFunc)(arg)
			}
		},
	},
	{
		funcName: "Caller5",
		funcDef:  Caller5,
		eval: func(t *testing.T) {
			if r := Caller5(); r != 0 {
				t.Fatalf("want result: %d, real: %d", 0, r)
			}
		},
		trampoline: func() interface{} {
			var result = func() int {
				return 60 + rand.Int()
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func() int {
				logger.LogTrace("proxy Caller5 called, no args")
				originFunc, _ := origin1.(*func() int)
				return (*originFunc)()
			}
		},
	},
	{
		funcName: "Caller6",
		funcDef:  Caller6,
		eval: func(t *testing.T) {
			if r := Caller6(3)(); r != 4 {
				t.Fatalf("want result: %d, real: %d", 4, r)
			}
		},
		trampoline: func() interface{} {
			var result = func(a int) func() int {
				return func() int {
					return a + 70
				}
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(a int) func() int {
				logger.LogTrace("proxy Caller6 called, args", a)
				originFunc, _ := origin1.(*func(a int) func() int)
				return (*originFunc)(a)
			}
		},
	},
	{
		funcName: "Caller7",
		funcDef:  Caller7,
		eval: func(t *testing.T) {
			Caller7(2)
		},
		trampoline: func() interface{} {
			var result = func(i int) {
				fmt.Println("trampoline")
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(i int) {
				logger.LogTrace("proxy Caller7 called, args", i)
				originFunc, _ := origin1.(*func(i int))
				(*originFunc)(i)
			}
		},
	},
	{
		funcName: "Caller8",
		funcDef:  Caller8,
		eval: func(t *testing.T) {
			if r := Caller8(-1); r != -1 {
				t.Fatalf("want result: %d, real: %d", -1, r)
			}
		},
		trampoline: func() interface{} {
			var result = func(i int) int {
				fmt.Println("trampoline")
				return 99
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(i int) func(int) int {
				logger.LogTrace("proxy Caller8 called, args", i)
				originFunc, _ := origin1.(*func(i int) func(int) int)
				return (*originFunc)(i)
			}
		},
	},
}

// main 测试静态代理
func TestProxy_fixIns(t *testing.T) {
	logger.LogLevel = logger.TraceLevel
	logger.SetLog2Console(true)
	for _, tc := range testCases1 {

		trampoline := tc.trampoline()

		// 静态代理函数
		patch, err := StaticProxyByName("git.code.oa.com/goom/mocker/internal/proxy."+
			tc.funcName, tc.proxy(trampoline), trampoline)
		if err != nil {
			t.Fatal("mock print err:", err)
		}

		patch.Apply()

		tc.eval(t)
		patch.Unpatch()
	}
	fmt.Println("all test is ok")
}
