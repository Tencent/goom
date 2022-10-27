// Package proxy 封装了给各种类型的代理(或叫 patch)中间层
// 负责比如外部传如私有函数名转换成 uintptr，trampoline 初始化，并发 proxy 等
package proxy

import (
	"reflect"
	"unsafe"

	"git.woa.com/goom/mocker/erro"
	"git.woa.com/goom/mocker/internal/hack"
	"git.woa.com/goom/mocker/internal/iface"
)

// Interface 构造接口代理，自动生成接口实现的桩指令织入到内存中
// ifaceVar 接口类型变量(指针类型)
// ctx 接口代理上下文
// method 代理模板方法名
// apply 代理函数, 代理函数的第一个参数类型必须是*IContext
// proxy 动态代理函数, 用于反射的方式回调, proxy 参数会覆盖 apply 参数值
// return error 异常
func Interface(ifaceVar interface{}, ctx *iface.IContext, method string, imp interface{}, proxy iface.PFunc) error {
	interfaceType := reflect.TypeOf(ifaceVar)
	if interfaceType.Kind() != reflect.Ptr {
		return erro.NewIllegalParamTypeError("interface Var", interfaceType.String(), "ptr")
	}

	typ := interfaceType.Elem()
	if typ.Kind() != reflect.Interface {
		return erro.NewIllegalParamTypeError("interface Var", typ.String(), "interface")
	}

	// check args len match
	argLen := reflect.TypeOf(imp).NumIn()
	funcTabIndex := methodIndexOf(typ, method)
	maxLen := typ.Method(funcTabIndex).Type.NumIn()
	if maxLen >= argLen {
		cause := erro.NewArgsNotMatchError(imp, argLen, maxLen+1)
		return erro.NewIllegalParamCError("interface As()", reflect.ValueOf(imp).String(), cause)
	}

	// 首次调用备份 iface
	gen := hack.UnpackEFace(ifaceVar).Data
	iface.BackUpTo(ctx, gen)

	// mock 接口方法
	var itabFunc = iface.GenCallableMethod(ctx, imp, proxy)
	// 上下文中查找接口代理对象的缓存
	ifaceCacheKey := typ.PkgPath() + "/" + typ.String()
	if fakeIface, ok := ctx.Cached(ifaceCacheKey); ok && !ctx.Canceled() {
		// 添加代理函数到 funcTab
		fakeIface.Tab.Fun[funcTabIndex] = itabFunc
		fakeIface.Data = unsafe.Pointer(ctx)
		applyIfaceTo(fakeIface, gen)
	} else {
		// 构造 iface 对象
		fakeIface = iface.MakeInterface(ctx, funcTabIndex, itabFunc, typ)
		ctx.Cache(ifaceCacheKey, fakeIface)
		applyIfaceTo(fakeIface, gen)
	}
	return nil
}

func methodIndexOf(typ reflect.Type, method string) int {
	funcTabIndex := 0
	// 根据方法名称获取到方法的 index
	for i := 0; i < typ.NumMethod(); i++ {
		if method == typ.Method(i).Name {
			funcTabIndex = i
			break
		}
	}
	return funcTabIndex
}

// applyIfaceTo 应用到变量
func applyIfaceTo(ifaceVar *hack.Iface, gen unsafe.Pointer) {
	// 伪造的 interface 赋值到指针变量
	*(*hack.Iface)(gen) = *ifaceVar
}
