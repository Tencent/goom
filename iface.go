// Package mocker 定义了mock的外层用户使用API定义,
// 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件实现了接口mock的能力。
package mocker

import (
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/proxy"
)

// InterfaceMocker 接口Mock
// 通过生成和替代接口变量实现Mock
type InterfaceMocker interface {
	ExportedMocker
	// Method 指定接口方法
	Method(name string) InterfaceMocker
	// As 将接口方法应用为函数类型
	// As调用之后,请使用Return或When API的方式来指定mock返回。
	// imp函数的第一个参数必须为*mocker.IContext, 作用是指定接口实现的接收体; 后续的参数原样照抄。
	As(imp interface{}) InterfaceMocker
	// Inject 将mock设置到变量
	Inject(iFace interface{}) InterfaceMocker
}

// IContext 接口mock的接收体
// 和internal/proxy.IContext保持同步
type IContext struct {
	// Data 可以传递任意数据
	Data interface{}
	// 占位属性
	_ unsafe.Pointer
}

// DefaultInterfaceMocker 默认接口Mocker
type DefaultInterfaceMocker struct {
	*baseMocker
	ctx     *proxy.IContext
	iFace   interface{}
	method  string
	funcDef interface{}
}

// NewDefaultInterfaceMocker 创建默认接口Mocker
// pkgName 包路径
// iFace 接口变量定义
func NewDefaultInterfaceMocker(pkgName string, iFace interface{}, ctx *proxy.IContext) *DefaultInterfaceMocker {
	return &DefaultInterfaceMocker{
		baseMocker: newBaseMocker(pkgName),
		ctx:        ctx,
		iFace:      iFace,
	}
}

// Method 指定mock的方法名
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

// Apply 应用接口方法mock为实际的接收体方法
// imp函数的第一个参数必须为*mocker.IContext, 作用是指定接口实现的接收体; 后续的参数原样照抄。
func (m *DefaultInterfaceMocker) Apply(imp interface{}) {
	if m.method == "" {
		panic("method is empty")
	}

	m.applyByIFaceMethod(m.ctx, m.iFace, m.method, imp, nil)
}

// As 将接口方法mock为实际的接收体方法
// imp函数的第一个参数必须为*mocker.IContext, 作用是指定接口实现的接收体; 后续的参数原样照抄。
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

// Origin 回调原函数(暂时不支持)
func (m *DefaultInterfaceMocker) Origin(interface{}) ExportedMocker {
	panic("implement me")
}

// Inject 回调原函数(暂时不支持)
func (m *DefaultInterfaceMocker) Inject(interface{}) InterfaceMocker {
	panic("implement me")
}

// If 回调原函数(暂时不支持)
func (m *DefaultInterfaceMocker) If() *If {
	panic("implement me")
}
