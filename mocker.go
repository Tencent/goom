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
	// 注意: Apply会覆盖之前设定的When条件和Return
	Apply(imp interface{})
	// Cancel 取消代理
	Cancel()
}

// ExportedMocker 导出函数mock接口
type ExportedMocker interface {
	Mocker
	// When 指定条件匹配
	When(args ...interface{}) *When
	// Matcher 执行返回值
	Return(args ...interface{}) *When
	// Origin 指定Mock之后的原函数, orign签名和mock的函数一致
	Origin(orign interface{}) ExportedMocker
	// If 条件表达式匹配
	If() *If
}

// UnexportedMocker 未导出函数mock接口
type UnexportedMocker interface {
	Mocker
	// As 将未导出函数(或方法)转换为导出函数(或方法)
	As(funcdef interface{}) ExportedMocker
	// Origin 指定Mock之后的原函数, orign签名和mock的函数一致
	Origin(orign interface{}) UnexportedMocker
}

// baseMocker mocker基础类型
type baseMocker struct {
	pkgName string
	origin  interface{}
	guard   *patch.PatchGuard
	imp     interface{}

	when *When
	_if  *If
}

func newBaseMocker(pkgName string) *baseMocker {
	return &baseMocker{
		pkgName: pkgName,
	}
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

// when 指定的返回值
func (m *baseMocker) whens(when *When) error {
	m.imp = reflect.MakeFunc(when.funcTyp, m.callback).Interface()
	m.when = when

	return nil
}

// if 指定的返回值
// func (m *baseMocker) ifs(_if *If) error {
// 	m.imp = reflect.MakeFunc(_if.funcTyp, m.callback).Interface()
// 	m._if = _if
//
// 	return nil
// }

func (m *baseMocker) callback(args []reflect.Value) (results []reflect.Value) {
	if m.when != nil {
		results = m.when.invoke(args)

		if results != nil {
			return results
		}
	}

	if m._if != nil {
		results = m._if.invoke(args)
		if results != nil {
			return results
		}
	}

	panic("not match any args, please spec default return use: mocker.Return()")
}

// Cancel 取消Mock
func (m *baseMocker) Cancel() {
	if m.guard != nil {
		m.guard.UnpatchWithLock()
	}

	m.when = nil
	m._if = nil
	m.origin = nil
}

// MethodMocker 对结构体函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type MethodMocker struct {
	*baseMocker
	structDef interface{}
	method    string
	methodIns interface{}
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

// Method 设置结构体的方法名
func (m *MethodMocker) Method(name string) ExportedMocker {
	if name == "" {
		panic("method is empty")
	}

	m.method = name
	sTyp := reflect.TypeOf(m.structDef)

	method, ok := sTyp.MethodByName(m.method)
	if !ok {
		panic("method " + m.method + " not found on " + sTyp.String())
	}

	m.methodIns = method.Func.Interface()

	return m
}

// ExportMethod 导出私有方法
func (m *MethodMocker) ExportMethod(name string) UnexportedMocker {
	if name == "" {
		panic("method is empty")
	}

	// 转换结构体名
	structName := getTypeName(m.structDef)
	if strings.Contains(structName, "*") {
		structName = fmt.Sprintf("(%s)", structName)
	}

	return (&UnexportedMethodMocker{
		baseMocker: m.baseMocker,
		structName: structName,
		methodName: name,
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
	if m.method == "" {
		panic("method is empty")
	}

	if m.when != nil {
		return m.when.When(args...)
	}

	sTyp := reflect.TypeOf(m.structDef)
	methodIns, ok := sTyp.MethodByName(m.method)

	if !ok {
		panic("method " + m.method + " not found on " + sTyp.String())
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, methodIns.Func.Interface(), args, nil, true); err != nil {
		panic(err)
	}

	if err := m.whens(when); err != nil {
		panic(err)
	}

	m.Apply(m.imp)

	return when
}

// Matcher 代理方法返回
func (m *MethodMocker) Return(returns ...interface{}) *When {
	if m.method == "" {
		panic("method is empty")
	}

	if m.when != nil {
		return m.when.Return(returns...)
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, m.methodIns, nil, returns, true); err != nil {
		panic(err)
	}

	if err := m.whens(when); err != nil {
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

// If 条件子句
func (m *MethodMocker) If() *If {
	return nil
}

// UnexportedMethodMocker 对结构体函数或方法进行mock
// 能支持到未导出类型、未导出类型的方法的Mock
type UnexportedMethodMocker struct {
	*baseMocker
	structName string
	methodName string
}

func (m *UnexportedMethodMocker) getObjName() string {
	return fmt.Sprintf("%s.%s.%s", m.pkgName, m.structName, m.methodName)
}

// Method 设置结构体的方法名
func (m *UnexportedMethodMocker) Method(name string) UnexportedMocker {
	m.methodName = name

	return m
}

// Apply 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *UnexportedMethodMocker) Apply(imp interface{}) {
	name := m.getObjName()
	if name == "" {
		panic("method name is empty")
	}

	if !strings.Contains(name, "*") {
		_, _ = unexports.FindFuncByName(name)
	}

	m.applyByName(name, imp)
}

// Origin 调用原函数
func (m *UnexportedMethodMocker) Origin(orign interface{}) UnexportedMocker {
	m.origin = orign

	return m
}

// As 将未导出函数(或方法)转换为导出函数(或方法)
func (m *UnexportedMethodMocker) As(funcdef interface{}) ExportedMocker {
	name := m.getObjName()
	if name == "" {
		panic("method name is empty")
	}

	var (
		err           error
		originFuncPtr uintptr
	)

	originFuncPtr, err = unexports.FindFuncByName(name)
	if err != nil {
		panic(err)
	}

	newFunc := unexports.NewFuncWithCodePtr(reflect.TypeOf(funcdef), originFuncPtr)

	return &DefMocker{
		baseMocker: m.baseMocker,
		funcdef:    newFunc.Interface(),
	}
}

// UnexportedFuncMocker 对函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type UnexportedFuncMocker struct {
	*baseMocker
	funcName string
}

func (m *UnexportedFuncMocker) getObjName() string {
	return fmt.Sprintf("%s.%s", m.pkgName, m.funcName)
}

// Apply 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *UnexportedFuncMocker) Apply(imp interface{}) {
	name := m.getObjName()

	m.applyByName(name, imp)
}

// Origin 调用原函数
func (m *UnexportedFuncMocker) Origin(orign interface{}) UnexportedMocker {
	m.origin = orign

	return m
}

// As 将未导出函数(或方法)转换为导出函数(或方法)
func (m *UnexportedFuncMocker) As(funcdef interface{}) ExportedMocker {
	name := m.getObjName()

	originFuncPtr, err := unexports.FindFuncByName(name)
	if err != nil {
		panic(err)
	}

	newFunc := unexports.NewFuncWithCodePtr(reflect.TypeOf(funcdef), originFuncPtr)

	return &DefMocker{
		baseMocker: m.baseMocker,
		funcdef:    newFunc.Interface(),
	}
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
	if m.when != nil {
		return m.when.When(args...)
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, m.funcdef, args, nil, false); err != nil {
		panic(err)
	}

	if err := m.whens(when); err != nil {
		panic(err)
	}

	m.Apply(m.imp)

	return when
}

// Matcher 代理方法返回
func (m *DefMocker) Return(returns ...interface{}) *When {
	if m.when != nil {
		return m.when.Return(returns...)
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, m.funcdef, nil, returns, false); err != nil {
		panic(err)
	}

	if err := m.whens(when); err != nil {
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

// If 条件子句
func (m *DefMocker) If() *If {
	return nil
}
