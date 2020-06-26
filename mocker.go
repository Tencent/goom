package mocker

import (
	"fmt"
	"strings"

	"git.code.oa.com/goom/mocker/internal/unexports"

	"git.code.oa.com/goom/mocker/internal/patch"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// Mocker mock接口
type Mocker interface {
	// Apply 代理方法实现
	Apply(imp interface{})
	// Return 代理方法返回
	Return(args ...interface{})
	// Cancel 取消代理
	Cancel()
}

// baseMocker mocker基础类型
type baseMocker struct {
	origin interface{}
	guard *patch.PatchGuard
}

// applyByName 根据函数名称应用mock
func (m *baseMocker) applyByName(funcname string, imp interface{}) (err error) {
	m.guard, err = proxy.StaticProxyByName(funcname, imp, m.origin)
	if err != nil {
		panic(fmt.Sprintf("proxy func definition error: %v", err))
	}
	m.guard.Apply()
	return
}

// applyByFunc 根据函数应用mock
func (m *baseMocker) applyByFunc(funcdef interface{}, imp interface{}) (err error) {
	m.guard, err = proxy.StaticProxyByFunc(funcdef, imp, m.origin)
	if err != nil {
		panic(fmt.Sprintf("proxy func definition error: %v", err))
	}
	m.guard.Apply()
	return
}

// Cancel 取消Mock
func (m *baseMocker) Cancel() {
	if m.guard != nil {
		m.guard.UnpatchWithLock()
	}
}

// MethodMocker 对结构体函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type MethodMocker struct {
	*baseMocker
	name   string
	namep  string
}

// Method 设置结构体的方法名
func (m *MethodMocker) Method(name string) Mocker {
	m.name = fmt.Sprintf("%s.%s", m.name, name)
	m.namep = fmt.Sprintf("%s.%s", m.namep, name)
	return m
}

// Apply 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *MethodMocker) Apply(imp interface{}) {
	if m.name == "" && m.namep == "" {
		panic("method name is empty")
	}

	var err error
	var mname = m.name
	_, err = unexports.FindFuncByName(m.name)
	if err != nil {
		mname = m.namep
	}

	m.applyByName(mname, imp)
}

// Return 代理方法返回
func (m *MethodMocker) Return(args ...interface{}) {
	panic("not implements")
}


// FuncMocker 对函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type FuncMocker struct {
	*baseMocker
	name   string
}

// Apply 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *FuncMocker) Apply(imp interface{}) {
	if m.name == "" {
		panic("func name is empty")
	}
	m.applyByName(m.name, imp)
}

// Return 代理方法返回
func (m *FuncMocker) Return(args ...interface{}) {
	panic("not implements")
}



// DefMocker 对函数或方法进行mock，使用函数定义筛选
type DefMocker struct {
	*baseMocker
	funcdef interface{}

	guard *patch.PatchGuard
}

// Apply 代理方法实现
func (m *DefMocker) Apply(imp interface{}) {
	if m.funcdef == nil {
		panic("funcdef is empty")
	}

	var funcname = getFunctionName(m.funcdef)
	if strings.HasSuffix(funcname, "-fm") {
		m.applyByName(strings.TrimRight(funcname, "-fm"), imp)
	} else {
		m.applyByFunc(m.funcdef, imp)
	}
}

// Return 代理方法返回
func (m *DefMocker) Return(args ...interface{}) {
	panic("not implements")
}