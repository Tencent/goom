package iface

import (
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/bytecode/stub"
	"git.code.oa.com/goom/mocker/internal/hack"
)

// IContext 接口 Mock 代码函数的接收体
// 避免被 mock 的接口变量为 nil, 无法通过单测逻辑中 mock==nil 的判断
type IContext struct {
	// Data 可以传递任意数据
	Data interface{}
	// 代理上下文数据
	p *PContext
}

// Cancel 取消接口代理
func (c *IContext) Cancel() {
	*c.p.originIface = *c.p.originIfaceValue
	c.p.canceled = true
}

// Canceled 是否已经被取消
func (c *IContext) Canceled() bool {
	return c.p.canceled
}

// Cached 获取缓存数据
func (c *IContext) Cached(key string) (v *hack.Iface, ok bool) {
	v, ok = c.p.ifaceCache[key]
	return
}

// Cache 缓存数据
func (c *IContext) Cache(key string, value *hack.Iface) {
	c.p.ifaceCache[key] = value
}

// NewContext 构造上下文
func NewContext() *IContext {
	return &IContext{
		Data: nil,
		p: &PContext{
			ifaceCache: make(map[string]*hack.Iface, 32),
		},
	}
}

// PContext 代理上下文
// 适配 proxy 包的 Context
type PContext struct {
	// ifaceCache iface 缓存
	ifaceCache map[string]*hack.Iface
	// originIface 原始接口地址
	originIface *hack.Iface
	// originIfaceValue 原始接口值
	originIfaceValue *hack.Iface
	// proxyFunc 代理函数, 需要内存持续持有
	proxyFunc reflect.Value
	// canceled 是否已经被取消
	canceled bool
}

// PFunc 代理函数类型的签名
type PFunc func(args []reflect.Value) (results []reflect.Value)

// notImplement 未实现的接口方法被调用的函数, 未配置 mock 的接口方法默认会跳转到调用此函数
func notImplement() {
	panic("method not implements. (please write a mocker on it)")
}

// MakeInterface 构造 interface 对象, 包含 receive、funcTab 等数据
func MakeInterface(ctx *IContext, funcTabIndex int, itabFunc uintptr, typ reflect.Type) *hack.Iface {
	funcTabData := [hack.MaxMethod]uintptr{}
	notImplements := reflect.ValueOf(notImplement).Pointer()
	for i := 0; i < hack.MaxMethod; i++ {
		funcTabData[i] = notImplements
	}
	funcTabData[funcTabIndex] = itabFunc

	// 伪造 iface
	structType := reflect.TypeOf(&IContext{})
	return &hack.Iface{
		Tab: &hack.Itab{
			Inter: (*uintptr)((*hack.Iface)(unsafe.Pointer(&typ)).Data),
			Type:  (*uintptr)((*hack.Iface)(unsafe.Pointer(&structType)).Data),
			Fun:   funcTabData,
		},
		Data: unsafe.Pointer(ctx),
	}
}

// BackUpTo 备份缓存 iface 指针到 IContext 中
func BackUpTo(ctx *IContext, iface unsafe.Pointer) {
	if ctx.p.originIfaceValue == nil {
		ctx.p.originIface = (*hack.Iface)(iface)
		originIfaceValue := *(*hack.Iface)(iface)
		ctx.p.originIfaceValue = &originIfaceValue
	}
}

// GenCallableMethod 生成可以直接 CALL 的接口方法实现, 带上下文 (rdx)
func GenCallableMethod(ctx *IContext, apply interface{}, proxy PFunc) uintptr {
	var (
		methodCaller uintptr
		err          error
	)

	if proxy == nil {
		// 生成桩代码,rdx 寄存器还原
		applyValue := reflect.ValueOf(apply)
		mockFuncPtr := (*hack.Value)(unsafe.Pointer(&applyValue)).Ptr
		methodCaller, err = MakeMethodCaller(mockFuncPtr)
	} else {
		// 生成桩代码,rdx 寄存器还原, 生成的调用将跳转到 proxy 函数
		methodTyp := reflect.TypeOf(apply)
		mockFunc := reflect.MakeFunc(methodTyp, proxy)
		callStub := reflect.ValueOf(stub.MakeFuncStub).Pointer()
		mockFuncPtr := (*hack.Value)(unsafe.Pointer(&mockFunc)).Ptr
		methodCaller, err = MakeMethodCallerWithCtx(mockFuncPtr, callStub)
		ctx.p.proxyFunc = mockFunc
	}

	if err != nil {
		panic(err)
	}
	return methodCaller
}
