package mocker

import (
	"fmt"
	"reflect"
	"strings"

	"git.code.oa.com/goom/mocker/internal/unexports"

	"git.code.oa.com/goom/mocker/internal/patch"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// Mocker mock接口
type Mocker interface {
	// Apply 代理方法实现
	Apply(imp interface{})
	// Cancel 取消代理
	Cancel()
}

// ExportedMocker 导出函数mock接口
type ExportedMocker interface {
	Mocker
	// When 指定条件匹配
	When(args ...interface{}) *When
	// Return 执行返回值
	Return(args ...interface{}) *When
	// 指定Mock之后的原函数, orign签名和mock的函数一致
	Origin(orign interface{}) ExportedMocker
}

// UnexportedMocker 未导出函数mock接口
type UnexportedMocker interface {
	Mocker
	// 指定Mock之后的原函数, orign签名和mock的函数一致
	Origin(orign interface{}) UnexportedMocker
}

// baseMocker mocker基础类型
type baseMocker struct {
	origin interface{}
	guard  *patch.PatchGuard
	imp    interface{}
}

// applyByName 根据函数名称应用mock
func (m *baseMocker) applyByName(funcname string, imp interface{}) {
	var err error

	m.guard, err = proxy.StaticProxyByName(funcname, imp, m.origin)
	if err != nil {
		panic(fmt.Sprintf("proxy func name error: %v", err))
	}

	m.guard.Apply()
	m.imp = imp
}

// applyByFunc 根据函数应用mock
func (m *baseMocker) applyByFunc(funcdef interface{}, imp interface{}) {
	var err error

	m.guard, err = proxy.StaticProxyByFunc(funcdef, imp, m.origin)
	if err != nil {
		panic(fmt.Sprintf("proxy func definition error: %v", err))
	}

	m.guard.Apply()
	m.imp = imp
}

// applyByMethod 根据函数名应用mock
func (m *MethodMocker) applyByMethod(structDef interface{}, method string, imp interface{}) {
	var err error

	m.guard, err = proxy.StaticProxyByMethod(reflect.TypeOf(structDef), method, imp, m.origin)
	if err != nil {
		panic(fmt.Sprintf("proxy method error: %v", err))
	}

	m.guard.Apply()
	m.imp = imp
}

// returns 指定的返回值
func (m *baseMocker) returns(when *When) error {
	m.imp = reflect.MakeFunc(when.funTyp, when.invoke).Interface()

	return nil
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
	pkgname   string
	structDef interface{}
	method    string
}

// Method 设置结构体的方法名
func (m *MethodMocker) Method(name string) ExportedMocker {
	m.method = name

	return m
}

// Method 设置结构体的方法名
func (m *MethodMocker) UnexportedMethod(name string) UnexportedMocker {
	return (&UnexportedMethodMocker{
		baseMocker: m.baseMocker,
		name:       fmt.Sprintf("%s.%s", m.pkgname, getTypeName(m.structDef)),
		namep:      fmt.Sprintf("%s.(%s)", m.pkgname, getTypeName(m.structDef)),
	}).Method(name)
}

// Apply 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *MethodMocker) Apply(imp interface{}) {
	if m.method == "" {
		panic("method is empty")
	}

	m.applyByMethod(m.structDef, m.method, imp)
}

// When 指定条件匹配
func (m *MethodMocker) When(args ...interface{}) *When {
	sTyp := reflect.TypeOf(m.structDef)

	methodIns, ok := sTyp.MethodByName(m.method)
	if !ok {
		panic("method " + m.method + " not found on " + sTyp.String())
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, methodIns.Func.Interface(), args, nil); err != nil {
		panic(err)
	}

	if err := m.returns(when); err != nil {
		panic(err)
	}

	m.Apply(m.imp)

	return when
}

// Return 代理方法返回
func (m *MethodMocker) Return(returns ...interface{}) *When {
	sTyp := reflect.TypeOf(m.structDef)

	methodIns, ok := sTyp.MethodByName(m.method)
	if !ok {
		panic("method " + m.method + " not found on " + sTyp.String())
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, methodIns.Func.Interface(), nil, returns); err != nil {
		panic(err)
	}

	if err := m.returns(when); err != nil {
		panic(err)
	}

	m.Apply(m.imp)

	return when
}

// Origin 调用原函数
func (m *MethodMocker) Origin(orign interface{}) ExportedMocker {
	m.origin = orign

	return m
}

// UnexportedMethodMocker 对结构体函数或方法进行mock
// 能支持到未导出类型、未导出类型的方法的Mock
type UnexportedMethodMocker struct {
	*baseMocker
	name  string
	namep string
}

// Method 设置结构体的方法名
func (m *UnexportedMethodMocker) Method(name string) UnexportedMocker {
	m.name = fmt.Sprintf("%s.%s", m.name, name)
	m.namep = fmt.Sprintf("%s.%s", m.namep, name)

	return m
}

// Apply 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *UnexportedMethodMocker) Apply(imp interface{}) {
	if m.name == "" && m.namep == "" {
		panic("method name is empty")
	}

	var (
		err   error
		mname = m.name
	)

	_, err = unexports.FindFuncByName(m.name)
	if err != nil {
		mname = m.namep
	}

	m.applyByName(mname, imp)
}

// Origin 调用原函数
func (m *UnexportedMethodMocker) Origin(orign interface{}) UnexportedMocker {
	m.origin = orign

	return m
}

// UnexportedFuncMocker 对函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type UnexportedFuncMocker struct {
	*baseMocker
	name string
}

// Apply 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *UnexportedFuncMocker) Apply(imp interface{}) {
	if m.name == "" {
		panic("func name is empty")
	}

	m.applyByName(m.name, imp)
}

// Origin 调用原函数
func (m *UnexportedFuncMocker) Origin(orign interface{}) UnexportedMocker {
	m.origin = orign

	return m
}

// DefMocker 对函数或方法进行mock，使用函数定义筛选
type DefMocker struct {
	*baseMocker
	funcdef interface{}
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

// When 指定条件匹配
func (m *DefMocker) When(args ...interface{}) *When {
	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, m.funcdef, args, nil); err != nil {
		panic(err)
	}

	if err := m.returns(when); err != nil {
		panic(err)
	}

	m.Apply(m.imp)

	return when
}

// Return 代理方法返回
func (m *DefMocker) Return(returns ...interface{}) *When {
	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, m.funcdef, nil, returns); err != nil {
		panic(err)
	}

	if err := m.returns(when); err != nil {
		panic(err)
	}

	m.Apply(m.imp)

	return when
}

// Origin 调用原函数
func (m *DefMocker) Origin(orign interface{}) ExportedMocker {
	m.origin = orign

	return m
}
