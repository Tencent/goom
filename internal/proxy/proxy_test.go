// Package proxy_test 对 proxy 包的测试
package proxy_test

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"testing"

	"git.code.oa.com/goom/mocker/internal/proxy"

	"git.code.oa.com/goom/mocker/internal/logger"
)

// 测试用例数据
var basePath = CurrentPackage()

var testCases = []*TestCase{
	{
		funcName: "Caller",
		funcDef:  Caller,
		eval: func() {
			_ = Caller(1000)
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(i int) int)(5)
		},
		trampoline: func() interface{} {
			var result = func(i int) int {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
				return i + 10
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			return func(i int) int {
				logger.LogDebug("proxy Caller called, args", i)
				originFunc, _ := origin.(*func(i int) int)
				return (*originFunc)(i)
			}
		},
	},
	{
		funcName: "Caller1",
		funcDef:  Caller1,
		eval: func() {
			_ = Caller1(5)
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(i int) int)(5)
		},
		trampoline: func() interface{} {
			var result = func(i int) int {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
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
		eval: func() {
			_ = Caller2(5)
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(i int) int)(5)
		},
		trampoline: func() interface{} {
			var result = func(i int) int {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
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
		eval: func() {
			//var arg = make(map[string]int, 0)
			Caller3(Arg{
				field1: field1,
				field2: nil,
			})
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(arg Arg) int)(Arg{
				field1: field1,
				field2: nil,
			})
		},
		trampoline: func() interface{} {
			var result = func(arg Arg) int {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
				return 40
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(arg Arg) int {
				logger.LogTrace("proxy Caller3 called, args", arg)
				originFunc, _ := origin1.(*func(arg Arg) int)
				return (*originFunc)(arg)
			}
		},
	},
	{
		funcName: "Caller4",
		funcDef:  Caller4,
		eval:     Caller4Eval,
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(arg *Arg) int)(&Arg{
				field1: field1,
				field2: nil,
			})
		},
		trampoline: func() interface{} {
			var result = func(arg *Arg) int {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
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
		eval: func() {
			Caller5()
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func() int)()
		},
		trampoline: func() interface{} {
			var result = func() int {
				fmt.Println("trampoline1")
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
		eval: func() {
			Caller6(3)()
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(a int) func() int)(3)()
		},
		trampoline: func() interface{} {
			var result = func(a int) func() int {
				return func() int {
					fmt.Println("trampoline1")
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
		eval: func() {
			Caller7(2)
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(a int))(2)
		},
		trampoline: func() interface{} {
			var result = func(i int) {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
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
		eval: func() {
			j := Caller8(5).inner.j
			if j < 0 {
				fmt.Println(j)
			}
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(i int) *Result)(5)
		},
		trampoline: func() interface{} {
			var result = func(i int) *Result {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
				return &Result{
					i: i * 100,
					inner: &InnerResult{
						j: i * 2 * 100,
					},
					m: make(map[string]int, 2),
				}
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(i int) *Result {
				logger.LogTrace("proxy Caller8 called, args", i)
				originFunc, _ := origin1.(*func(i int) *Result)
				return (*originFunc)(i)
			}
		},
	},
	{
		funcName: "Caller9",
		funcDef:  Caller9,
		eval: func() {
			j := Caller9(5).m
			if len(j) > 0 {
				fmt.Println(j)
			}
		},
		evalMakeFunc: func(makeFunc interface{}) {
			makeFunc.(func(i int) Result)(5)
		},
		trampoline: func() interface{} {
			var result = func(i int) Result {
				fmt.Println("trampoline")
				fmt.Println("trampoline1")
				return Result{
					i: i * 100,
					inner: &InnerResult{
						j: i * 2 * 100,
					},
					m: make(map[string]int, 2),
				}
			}
			return &result
		},
		proxy: func(origin interface{}) interface{} {
			var origin1 = origin
			return func(i int) Result {
				logger.LogTrace("proxy Caller9 called, args", i)
				originFunc, _ := origin1.(*func(i int) Result)
				return (*originFunc)(i)
			}
		},
	},
}

// TestTestStaticProxy 测试静态代理
func TestTestStaticProxy(t *testing.T) {
	logger.LogLevel = logger.DebugLevel
	logger.SetLog2Console(true)

	for _, tc := range testCases {
		trampoline := tc.trampoline()

		// 静态代理函数
		patch, err := proxy.StaticProxyByName(basePath+"."+tc.funcName, tc.proxy(trampoline), trampoline)
		if err != nil {
			log.Println("mock print err:", err)
			continue
		}

		tc.eval()
		patch.Unpatch()
	}

	fmt.Println("ok")
}

// TestTestStaticProxy 测试静态代理
func BenchmarkStaticProxy(b *testing.B) {
	logger.LogLevel = logger.TraceLevel
	logger.SetLog2Console(true)

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			trampoline := tc.trampoline()

			fun := tc.proxy(trampoline)

			// 静态代理函数
			patch, err := proxy.StaticProxyByName(basePath+"."+tc.funcName, fun, trampoline)
			if err != nil {
				b.Errorf("mock %s print err:%s", tc.funcName, err)
			}

			tc.eval()
			patch.Unpatch()
		}
	}
}

// TestStaticProxyConcurrent 测试并发支持
func TestStaticProxyConcurrent(t *testing.T) {
	logger.LogLevel = logger.WarningLevel
	logger.SetLog2Console(true)

	wait := make(chan int)

	for c := 0; c < 10; c++ {
		go func(c1 int) {
			for i := 0; i < 100; i++ {
				for _, tc := range testCases {
					trampoline := tc.trampoline()

					// 静态代理函数
					patch, err := proxy.StaticProxyByName(basePath+"."+tc.funcName, tc.proxy(trampoline), trampoline)
					if err != nil {
						t.Error("mock print err:", err)
					}

					tc.eval()
					patch.Unpatch()

					wait <- i * c1
				}
			}
		}(c)
	}

	for i := 0; i < 10*100*len(testCases); i++ {
		<-wait
	}
}

// TestConcurrent 测试运行中 patch 并发支持
func TestStaticProxyConcurrent1(t *testing.T) {
	logger.LogLevel = logger.WarningLevel
	logger.SetLog2Console(true)

	for c := 0; c < 50; c++ {
		go func() {
			for i := 0; i < 10000; i++ {
				for _, t := range testCases {
					t.eval()
				}
			}
		}()
	}

	for c := 0; c < 1000; c++ {
		go func() {
			for _, tc := range testCases {
				trampoline := tc.trampoline()

				// 静态代理函数
				patch, err := proxy.StaticProxyByName(basePath+"."+tc.funcName, tc.proxy(trampoline), trampoline)
				if err != nil {
					t.Error("mock print err:", err)
				}

				tc.eval()
				patch.Unpatch()
			}
		}()
	}

	wait := make(chan int)

	for c := 0; c < 50; c++ {
		go func(c1 int) {
			for i := 0; i < 10000; i++ {
				for _, t := range testCases {
					t.eval()
					wait <- i * c1
				}
			}
		}(c)
	}

	for i := 0; i < 50*10000*len(testCases); i++ {
		<-wait
	}
}

// TestConcurrent 测试运行中 patch 并发支持
// TODO fix nil pointer
func TestStaticProxyConcurrentOnce(t *testing.T) {
	logger.LogLevel = logger.InfoLevel
	logger.SetLog2Console(true)

	for c := 0; c < 50; c++ {
		go func() {
			for i := 0; i < 10000; i++ {
				for _, t := range testCases {
					t.eval()
				}
			}
		}()
	}

	go func() {
		for _, tc := range testCases {
			trampoline := tc.trampoline()

			// 静态代理函数
			patch, err := proxy.StaticProxyByName(basePath+"."+tc.funcName, tc.proxy(trampoline), trampoline)
			if err != nil {
				t.Error("mock print err:", err)
			}

			tc.eval()
			patch.UnpatchWithLock()
		}
	}()

	wait := make(chan int)

	for c := 0; c < 50; c++ {
		go func(c1 int) {
			for i := 0; i < 100; i++ {
				for _, t := range testCases {
					t.eval()
					wait <- i * c1
				}
			}
		}(c)
	}

	for i := 0; i < 50*100*len(testCases); i++ {
		<-wait
	}
}

// CurrentPackage 获取当前调用的包路径
func CurrentPackage() string {
	return currentPackage(2)
}

// currentPackage 获取调用者的包路径
func currentPackage(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	callerName := runtime.FuncForPC(pc).Name()

	if i := strings.Index(callerName, ".("); i > -1 {
		return callerName[:i]
	}

	if i := strings.LastIndex(callerName, "."); i > -1 {
		return callerName[:i]
	}

	return callerName
}
