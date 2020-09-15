package proxy

import (
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/errortype"

	"git.code.oa.com/goom/mocker/internal/hack"
	"git.code.oa.com/goom/mocker/internal/stub"
)

// IContext 接口实现定义
type IContext struct {
	// Data 可以传递任意数据
	Data interface{}
	// 代理上下文数据
	p *PContext
}

// PContext 代理上下文
type PContext struct {
	// ifaceCache iface 缓存
	ifaceCache map[string]*hack.Iface
	// originIface 原始接口地址
	originIface *hack.Iface
	// originIfaceValue 原始接口值
	originIfaceValue *hack.Iface
	// 需要指针持有的变量
	// proxyfunc
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

// MakeInterfaceImpl 构造接口代理
// iface 接口类型变量,指针类型
// p 接口代理上下文
// method 代理模板方法名
// apply 代理函数, 代理函数的第一个参数类型必须是*IContext
// proxy 动态代理函数, 用于反射的方式回调, proxy参数会覆盖apply参数值
// return error 异常
func MakeInterfaceImpl(iface interface{}, ctx *IContext, method string,
	apply interface{}, proxy func(args []reflect.Value) (results []reflect.Value)) error {
	ifaceType := reflect.TypeOf(iface)
	if ifaceType.Kind() != reflect.Ptr {
		return errortype.NewIllegalParamTypeError("iface", ifaceType.String(), "ptr")
	}

	typ := ifaceType.Elem()
	funcTabIndex := 0

	for i := 0; i < typ.NumMethod(); i++ {
		if method == typ.Method(i).Name {
			funcTabIndex = i
			break
		}
	}

	gen := hack.UnpackEFace(iface).Data

	// mock接口方法
	var itabFunc = genCallableFunc(ctx, apply, proxy)

	ifaceCacheKey := typ.PkgPath() + "/" + typ.String()
	if iface, ok := ctx.p.ifaceCache[ifaceCacheKey]; ok {
		iface.Tab.Fun[funcTabIndex] = itabFunc
		if ctx != nil {
			iface.Data = unsafe.Pointer(ctx)
		}
	} else {

		funcTabData := [hack.MaxMethod]uintptr{}
		notImplements := reflect.ValueOf(notImplement).Pointer()
		for i := 0; i < hack.MaxMethod; i++ {
			funcTabData[i] = notImplements
		}
		funcTabData[funcTabIndex] = itabFunc

		iface := hack.Iface{
			Tab: &hack.Itab{
				Fun: funcTabData,
			},
			Data: unsafe.Pointer(ctx),
		}

		// 首次调用备份iface
		if ctx.p.originIfaceValue == nil {

			ctx.p.originIface = (*hack.Iface)(unsafe.Pointer(gen))

			originIfaceValue := *(*hack.Iface)(unsafe.Pointer(gen))
			ctx.p.originIfaceValue = &originIfaceValue
		}

		// 伪装iface
		*(*hack.Iface)(unsafe.Pointer(gen)) = iface

		ctx.p.ifaceCache[ifaceCacheKey] = &iface
	}

	return nil
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

		genStub, err = stub.GenStub(mockFuncPtr)
		if err != nil {
			panic(err)
		}
	} else {
		// 生成桩代码,rdx寄存器还原
		methodTyp := reflect.TypeOf(apply)
		mockfunc := reflect.MakeFunc(methodTyp, proxy)
		callStub := reflect.ValueOf(stub.MakeFuncStub).Pointer()

		mockFuncPtr := (*hack.Value)(unsafe.Pointer(&mockfunc)).Ptr

		genStub, err = stub.GenStubWithCtx(mockFuncPtr, callStub)
		if err != nil {
			panic(err)
		}

		ctx.p.proxyFunc = mockfunc
	}

	return genStub
}
