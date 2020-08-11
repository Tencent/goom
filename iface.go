package mocker

import (
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/proxy"
)

// InterfaceMocker 接口Mock
// 生成和替代接口变量实现Mock
type InterfaceMocker interface {
	Mocker
	// Method 指定接口方法
	Method(name string) InterfaceMocker
	// As 将接口方法应用为函数类型
	As(funcdef interface{}) InterfaceMocker
	// When 指定条件匹配
	When(args ...interface{}) *When
	// Matcher 执行返回值
	Return(args ...interface{}) *When
	// Origin 指定原接口变量, orign类型和mock的函数一致
	Origin(orign interface{}) InterfaceMocker
	// Inject 将mock设置到变量
	Inject(iface interface{}) InterfaceMocker
	// If 条件表达式匹配
	If() *If
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
	ctx    *proxy.IContext
	iface  interface{}
	method string
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

	sTyp := reflect.TypeOf(m.iface).Elem()

	_, ok := sTyp.MethodByName(name)
	if !ok {
		panic("method " + name + " not found on " + sTyp.String())
	}

	m.method = name

	return m
}

func (m *DefaultInterfaceMocker) Apply(imp interface{}) {
	if m.method == "" {
		panic("method is empty")
	}

	m.applyByIfaceMethod(m.ctx, m.iface, m.method, imp)
}

func (m *DefaultInterfaceMocker) As(funcdef interface{}) InterfaceMocker {
	panic("implement me")
}

func (m *DefaultInterfaceMocker) When(args ...interface{}) *When {
	panic("implement me")
}

func (m *DefaultInterfaceMocker) Return(args ...interface{}) *When {
	panic("implement me")
}

func (m *DefaultInterfaceMocker) Origin(orign interface{}) InterfaceMocker {
	panic("implement me")
}

func (m *DefaultInterfaceMocker) Inject(iface interface{}) InterfaceMocker {
	panic("implement me")
}

func (m *DefaultInterfaceMocker) If() *If {
	panic("implement me")
}
