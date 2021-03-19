// Packge proxy封装了给各种类型的代理(或较patch)中间层
// 负责比如外部传如私有函数名转换成uintptr，trampoline初始化，并发proxy等
package proxy

import (
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/errobj"

	"git.code.oa.com/goom/mocker/internal/hack"
	"git.code.oa.com/goom/mocker/internal/stub"
)

// IContext 接口Mock代码函数的接收体
// 避免被mock的接口变量为nil, 无法通过单测逻辑中mocki==nil的判断
type IContext struct {
	// Data 可以传递任意数据
	Data interface{}
	// 代理上下文数据
	p *PContext
}

// PContext 代理上下文
// 适配proxy包的Context
type PContext struct {
	// ifaceCache iface 缓存
	ifaceCache map[string]*hack.Iface
	// originIface 原始接口地址
	originIface *hack.Iface
	// originIfaceValue 原始接口值
	originIfaceValue *hack.Iface
	// proxyfunc 代理函数, 需要内存持续持有
	proxyFunc reflect.Value
}

// Cancel 取消接口代理
func (c *IContext) Cancel() {
	*c.p.originIface = *c.p.originIfaceValue
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

// notImplement 未实现的接口方法被调用的函数
func notImplement() {
	panic("method not implements. (please write a mocker on it)")
}

// MakeInterfaceImpl 构造接口代理，自动生成接口实现的桩指令织入到内存中
// iface 接口类型变量,指针类型
// ctx 接口代理上下文
// method 代理模板方法名
// apply 代理函数, 代理函数的第一个参数类型必须是*IContext
// proxy 动态代理函数, 用于反射的方式回调, proxy参数会覆盖apply参数值
// return error 异常
func MakeInterfaceImpl(iface interface{}, ctx *IContext, method string,
	apply interface{}, proxy func(args []reflect.Value) (results []reflect.Value)) error {
	ifaceType := reflect.TypeOf(iface)
	if ifaceType.Kind() != reflect.Ptr {
		return errobj.NewIllegalParamTypeError("iface", ifaceType.String(), "ptr")
	}

	typ := ifaceType.Elem()
	if typ.Kind() != reflect.Interface {
		return errobj.NewIllegalParamTypeError("iface var", typ.String(), "interface")
	}

	funcTabIndex := 0

	// 根据方法名称获取到方法的index
	for i := 0; i < typ.NumMethod(); i++ {
		if method == typ.Method(i).Name {
			funcTabIndex = i
			break
		}
	}

	gen := hack.UnpackEFace(iface).Data

	// 首次调用备份iface
	backUp2Context(ctx, gen)

	// mock接口方法
	var itabFunc = genCallableFunc(ctx, apply, proxy)

	ifaceCacheKey := typ.PkgPath() + "/" + typ.String()
	// 上下文中查找接口代理对象的缓存
	if iface, ok := ctx.p.ifaceCache[ifaceCacheKey]; ok {
		iface.Tab.Fun[funcTabIndex] = itabFunc
		if ctx != nil {
			iface.Data = unsafe.Pointer(ctx)
		}

		return nil
	}

	// 构造funcTab数据
	funcTabData := [hack.MaxMethod]uintptr{}
	notImplements := reflect.ValueOf(notImplement).Pointer()
	for i := 0; i < hack.MaxMethod; i++ {
		funcTabData[i] = notImplements
	}
	funcTabData[funcTabIndex] = itabFunc

	// 伪造iface
	structType := reflect.TypeOf(&IContext{})
	fakeIface := hack.Iface{
		Tab: &hack.Itab{
			Inter: (*uintptr)((*hack.Iface)(unsafe.Pointer(&typ)).Data),
			Type:  (*uintptr)((*hack.Iface)(unsafe.Pointer(&structType)).Data),
			Fun:   funcTabData,
		},
		Data: unsafe.Pointer(ctx),
	}

	// 伪造的iface赋值到指针变量
	*(*hack.Iface)(unsafe.Pointer(gen)) = fakeIface

	ctx.p.ifaceCache[ifaceCacheKey] = &fakeIface

	return nil
}

// backUp2Context 备份缓存iface指针到IContext中
func backUp2Context(ctx *IContext, iface unsafe.Pointer) {
	if ctx.p.originIfaceValue == nil {

		ctx.p.originIface = (*hack.Iface)(unsafe.Pointer(iface))

		originIfaceValue := *(*hack.Iface)(unsafe.Pointer(iface))
		ctx.p.originIfaceValue = &originIfaceValue
	}
}

// genCallableFunc 生成可以直接CALL的函数, 带上下文(rdx)
func genCallableFunc(ctx *IContext, apply interface{},
	proxy func(args []reflect.Value) (results []reflect.Value)) uintptr {
	var (
		genStub uintptr
		err     error
	)

	if proxy == nil {
		// 生成桩代码,rdx寄存去还原
		applyValue := reflect.ValueOf(apply)
		mockFuncPtr := (*hack.Value)(unsafe.Pointer(&applyValue)).Ptr

		genStub, err = stub.MakeIfaceCaller(mockFuncPtr)
		if err != nil {
			panic(err)
		}
	} else {
		// 生成桩代码,rdx寄存器还原
		methodTyp := reflect.TypeOf(apply)
		mockfunc := reflect.MakeFunc(methodTyp, proxy)
		callStub := reflect.ValueOf(stub.MakeFuncStub).Pointer()

		mockFuncPtr := (*hack.Value)(unsafe.Pointer(&mockfunc)).Ptr

		genStub, err = stub.MakeIfaceCallerWithCtx(mockFuncPtr, callStub)
		if err != nil {
			panic(err)
		}

		ctx.p.proxyFunc = mockfunc
	}

	return genStub
}
