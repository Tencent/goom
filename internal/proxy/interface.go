package proxy

import (
	"git.code.oa.com/goom/mocker/internal/hack"
	"git.code.oa.com/goom/mocker/internal/stub"
	"reflect"
	"unsafe"
)

// IContext 接口实现定义
type IContext struct {
	// Data 可以传递任意数据
	Data interface{}
}

var ifaceCache = make(map[string]*hack.Iface, 64)

// MakeInterfaceImpl 构造接口代理
func MakeInterfaceImpl(iface interface{}, ctx *IContext, method string,
		apply interface{}, proxy func(args []reflect.Value) (results []reflect.Value)) error {
	typ := reflect.TypeOf(iface).Elem()

	funcTabIndex := 0
	for i := 0; i < typ.NumMethod(); i++ {
		if method == typ.Method(i).Name {
			funcTabIndex = i
			break
		}
	}

	gen := hack.UnpackEFace(iface).Data

	// mock接口方法
	var itabFunc uintptr
	if proxy == nil {
		applyValue := reflect.ValueOf(apply)
		mockFuncPtr := (*hack.Value)(unsafe.Pointer(&applyValue)).Ptr
		genStub, err := stub.GenStub(mockFuncPtr)
		if err != nil {
			panic(err)
		}
		itabFunc = genStub
	} else {
		methodTyp := reflect.TypeOf(apply)
		mockfunc := reflect.MakeFunc(methodTyp, proxy)
		callStub := reflect.ValueOf(stub.MakeFuncStub).Pointer()

		mockFuncPtr := (*hack.Value)(unsafe.Pointer(&mockfunc)).Ptr
		genStub, err := stub.GenStubWithCtx(mockFuncPtr, callStub)
		if err != nil {
			panic(err)
		}
		itabFunc = genStub
	}

	ifaceCacheKey := typ.PkgPath() + "/" + typ.String()
	if iface, ok := ifaceCache[ifaceCacheKey]; ok {
		iface.Tab.Fun[funcTabIndex] = itabFunc
		if ctx != nil {
			iface.Data = unsafe.Pointer(ctx)
		}
	} else {

		funcTabData := [99]uintptr{}
		funcTabData[funcTabIndex] = itabFunc

		iface := hack.Iface{
			Tab: &hack.Itab{
				Fun: funcTabData,
			},
			Data: unsafe.Pointer(ctx),
		}

		// 伪装iface
		*(*hack.Iface)(unsafe.Pointer(gen)) = iface

		ifaceCache[ifaceCacheKey] = &iface
	}
	return nil
}