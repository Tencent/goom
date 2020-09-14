package mocker

import (
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/proxy"
)

// InterfaceMocker 接口Mock
// 生成和替代接口变量实现Mock
type InterfaceMocker interface {
	ExportedMocker
	// Method 指定接口方法
	Method(name string) InterfaceMocker
	// As 将接口方法应用为函数类型
	As(funcdef interface{}) InterfaceMocker
	// Inject 将mock设置到变量
	Inject(iface interface{}) InterfaceMocker
	// If 条件表达式匹配
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
	iface   interface{}
	method  string
	funcDef interface{}
}

// NewDefaultInterfaceMocker 创建默认接口Mocker
// pkgName 包路径
// iface 接口变量定义
func NewDefaultInterfaceMocker(pkgName string, iface interface{}, ctx *proxy.IContext) *DefaultInterfaceMocker {
	return &DefaultInterfaceMocker{
		baseMocker: newBaseMocker(pkgName),
		ctx:        ctx,
		iface:      iface,
	}
}

func (m *DefaultInterfaceMocker) Method(name string) InterfaceMocker {
	if name == "" {
		panic("method is empty")
	}

	m.checkMethod(name)
	m.method = name

	return m
}

func (m *DefaultInterfaceMocker) checkMethod(name string) {
	sTyp := reflect.TypeOf(m.iface).Elem()

	_, ok := sTyp.MethodByName(name)
	if !ok {
		panic("method " + name + " not found on " + sTyp.String())
	}
}

func (m *DefaultInterfaceMocker) Apply(imp interface{}) {
	if m.method == "" {
		panic("method is empty")
	}

	m.applyByIfaceMethod(m.ctx, m.iface, m.method, imp, nil)
}

func (m *DefaultInterfaceMocker) As(funcdef interface{}) InterfaceMocker {
	if m.method == "" {
		panic("method is empty")
	}

	m.funcDef = funcdef

	return m
}

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

	m.applyByIfaceMethod(m.ctx, m.iface, m.method, m.funcDef, m.callback)
	m.when = when

	return when
}

func (m *DefaultInterfaceMocker) Return(returns ...interface{}) *When {
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

	m.applyByIfaceMethod(m.ctx, m.iface, m.method, m.funcDef, m.callback)
	m.when = when

	return when
}

func (m *DefaultInterfaceMocker) Origin(orign interface{}) ExportedMocker {
	panic("implement me")
}

func (m *DefaultInterfaceMocker) Inject(iface interface{}) InterfaceMocker {
	panic("implement me")
}

func (m *DefaultInterfaceMocker) If() *If {
	panic("implement me")
}