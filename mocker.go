package mocker

import (
	"fmt"
	"git.code.oa.com/goom/mocker/internal/unexports"
	"reflect"
	"runtime"
	"strings"

	"git.code.oa.com/goom/mocker/internal/patch"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// Mocker mock接口
type Mocker interface {
	// Proxy 代理方法实现
	Proxy(imp interface{})
	// Return 代理方法返回
	Return()
	// Cancel 取消代理
	Cancel()
}

// MethodMocker 对结构体函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type MethodMocker struct {
	name   string
	namep  string
	origin interface{}

	guard *patch.PatchGuard
}

// Method 设置结构体的方法名
func (m *MethodMocker) Method(name string) Mocker {
	m.name = fmt.Sprintf("%s.%s", m.name, name)
	m.namep = fmt.Sprintf("%s.%s", m.namep, name)

	return m
}

// Proxy 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *MethodMocker) Proxy(imp interface{}) {
	if m.name == "" && m.namep == "" {
		panic("method name is empty")
	}

	var err error
	var mname = m.name
	_, err = unexports.FindFuncByName(m.name)
	if err != nil {
		mname = m.namep
	}

	m.guard, err = proxy.StaticProxyByName(mname, imp, m.origin)
	if err != nil {
		panic(fmt.Sprintf("proxy method error: %v", err))
	}

	m.guard.Apply()
}

// Return 代理方法返回
func (m *MethodMocker) Return() {
	panic("not implements")
}

// Cancel 取消Mock
func (m *MethodMocker) Cancel() {
	if m.guard != nil {
		m.guard.UnpatchWithLock()
	}
}

// FuncMocker 对函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type FuncMocker struct {
	name   string
	origin interface{}

	guard *patch.PatchGuard
}

// Proxy 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *FuncMocker) Proxy(imp interface{}) {
	if m.name == "" {
		panic("func name is empty")
	}

	var err error

	m.guard, err = proxy.StaticProxyByName(m.name, imp, m.origin)
	if err != nil {
		panic(fmt.Sprintf("proxy func error: %v", err))
	}

	m.guard.Apply()
}

// Return 代理方法返回
func (m *FuncMocker) Return() {
	panic("not implements")
}

// Cancel 取消Mock
func (m *FuncMocker) Cancel() {
	if m.guard != nil {
		m.guard.UnpatchWithLock()
	}
}

// DefMocker 对函数或方法进行mock，使用函数定义筛选
type DefMocker struct {
	funcdef interface{}
	origin  interface{}

	guard *patch.PatchGuard
}

// Proxy 代理方法实现
func (m *DefMocker) Proxy(imp interface{}) {
	if m.funcdef == nil {
		panic("funcdef is empty")
	}

	var err error
	var funcname = getFunctionName(m.funcdef)
	if strings.HasSuffix(funcname, "-fm") {
		m.guard, err = proxy.StaticProxyByName(strings.TrimRight(funcname, "-fm"), imp, m.origin)
		if err != nil {
			panic(fmt.Sprintf("proxy func definition error: %v", err))
		}
	} else {
		m.guard, err = proxy.StaticProxyByFunc(m.funcdef, imp, m.origin)
		if err != nil {
			panic(fmt.Sprintf("proxy func definition error: %v", err))
		}
	}

	m.guard.Apply()
}

// Return 代理方法返回
func (m *DefMocker) Return() {
	panic("not implements")
}

// Cancel 取消Mock
func (m *DefMocker) Cancel() {
	if m.guard != nil {
		m.guard.UnpatchWithLock()
	}
}

// getFunctionName 获取函数名称
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}