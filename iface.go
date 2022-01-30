// Package mocker 定义了 mock 的外层用户使用 API 定义,
// 包括函数、方法、接口、未导出函数(或方法的)的 Mocker 的实现。
// 当前文件实现了接口 mock 的能力。
package mocker

import (
	"fmt"
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// IContext 接口 mock 的接收体
// 和 internal/proxy.IContext 保持同步
type IContext struct {
	// Data 可以传递任意数据
	Data interface{}
	// 占位属性
	_ unsafe.Pointer
}

// InterfaceMocker 接口 Mock
// 通过生成和替代接口变量实现 Mock
type InterfaceMocker interface {
	ExportedMocker
	// Method 指定接口方法
	Method(name string) InterfaceMocker
	// As 将接口方法应用为函数类型
	// As 调用之后,请使用 Return 或 When API 的方式来指定 mock 返回。
	// imp 函数的第一个参数必须为*mocker.IContext, 作用是指定接口实现的接收体; 后续的参数原样照抄。
	As(imp interface{}) InterfaceMocker
	// Inject 将 mock 设置到变量
	Inject(iFace interface{}) InterfaceMocker
}

// DefaultInterfaceMocker 默认接口 Mocker
type DefaultInterfaceMocker struct {
	*baseMocker
	ctx     *proxy.IContext
	iFace   interface{}
	method  string
	funcDef interface{}
}

// String 接口 Mock 名称
func (m *DefaultInterfaceMocker) String() string {
	t := reflect.TypeOf(m.iFace)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return fmt.Sprintf("%s.%s", t.String(), m.method)
}

// NewDefaultInterfaceMocker 创建默认接口 Mocker
// pkgName 包路径
// iFace 接口变量定义
func NewDefaultInterfaceMocker(pkgName string, iFace interface{}, ctx *proxy.IContext) *DefaultInterfaceMocker {
	return &DefaultInterfaceMocker{
		baseMocker: newBaseMocker(pkgName),
		ctx:        ctx,
		iFace:      iFace,
	}
}

// Method 指定 mock 的方法名
func (m *DefaultInterfaceMocker) Method(name string) InterfaceMocker {
	if name == "" {
		panic("method is empty")
	}

	m.checkMethod(name)
	m.method = name

	return m
}

// checkMethod 检查是否能找到函数
func (m *DefaultInterfaceMocker) checkMethod(name string) {
	sTyp := reflect.TypeOf(m.iFace).Elem()

	_, ok := sTyp.MethodByName(name)
	if !ok {
		panic("method " + name + " not found on " + sTyp.String())
	}
}

// Apply 应用接口方法 mock 为实际的接收体方法
// imp 函数的第一个参数必须为*mocker.IContext, 作用是指定接口实现的接收体; 后续的参数原样照抄。
func (m *DefaultInterfaceMocker) Apply(imp interface{}) {
	if m.method == "" {
		panic("method is empty")
	}

	m.applyByIFaceMethod(m.ctx, m.iFace, m.method, imp, nil)
}

// As 将接口方法 mock 为实际的接收体方法
// imp 函数的第一个参数必须为*mocker.IContext, 作用是指定接口实现的接收体; 后续的参数原样照抄。
func (m *DefaultInterfaceMocker) As(imp interface{}) InterfaceMocker {
	if m.method == "" {
		panic("method is empty")
	}

	m.funcDef = imp

	return m
}

// When 执行参数匹配时的返回值
func (m *DefaultInterfaceMocker) When(args ...interface{}) *When {
	if m.method == "" {
		panic("method is empty")
	}

	if m.when != nil {
		return m.when.When(args...)
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, m.funcDef, args, nil, true); err != nil {
		panic(err)
	}

	m.applyByIFaceMethod(m.ctx, m.iFace, m.method, m.funcDef, m.callback)
	m.when = when

	return when
}

// Return 指定返回值
func (m *DefaultInterfaceMocker) Return(returns ...interface{}) *When {
	if m.funcDef == nil {
		panic("must use As() API before call Return()")
	}

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

	if when, err = CreateWhen(m, m.funcDef, nil, returns, true); err != nil {
		panic(err)
	}

	m.applyByIFaceMethod(m.ctx, m.iFace, m.method, m.funcDef, m.callback)
	m.when = when

	return when
}

// Returns 指定返回多个值
func (m *DefaultInterfaceMocker) Returns(returns ...interface{}) *When {
	if m.funcDef == nil {
		panic("must use As() API before call Return()")
	}

	if m.method == "" {
		panic("method is empty")
	}

	if m.when != nil {
		return m.when.Returns(returns...)
	}

	var (
		when *When
		err  error
	)

	if when, err = CreateWhen(m, m.funcDef, nil, nil, true); err != nil {
		panic(err)
	}

	when.Returns(returns...)
	m.applyByIFaceMethod(m.ctx, m.iFace, m.method, m.funcDef, m.callback)
	m.when = when

	return when
}

// Origin 回调原函数(暂时不支持)
func (m *DefaultInterfaceMocker) Origin(interface{}) ExportedMocker {
	panic("implement me")
}

// Inject 回调原函数(暂时不支持)
func (m *DefaultInterfaceMocker) Inject(interface{}) InterfaceMocker {
	panic("implement me")
}

// applyByIFaceMethod 根据接口方法应用 mock
func (m *DefaultInterfaceMocker) applyByIFaceMethod(ctx *proxy.IContext, iFace interface{},
	method string, imp interface{}, implV proxy.PFunc) {
	imp, implV = interceptDebugInfo(imp, implV, m)
	m.baseMocker.applyByIFaceMethod(ctx, iFace, method, imp, implV)

	logger.Log2Consolefc(logger.DebugLevel, "mocker [%s] apply.", logger.Caller(6), m.String())
}
