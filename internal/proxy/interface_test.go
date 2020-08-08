package proxy_test

import (
	"fmt"
	"git.code.oa.com/goom/mocker/internal/proxy"
	"reflect"
	"runtime/debug"
	"testing"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/stub"

	"git.code.oa.com/goom/mocker/internal/hack"
)

// TestInterfaceCall 测试接口调用
func TestInterfaceCall(t *testing.T) {
	i := getImpl(1)
	i.Call(0)
}

// TestAutoGen 测试生成任意接口实现
func TestAutoGen(t *testing.T) {
	logger.LogLevel = logger.DebugLevel
	logger.Log2Console(true)

	gen := (I)(nil)

	proxy.MakeInterfaceImpl(&gen, &proxy.IContext{}, "Call", func(ctx *proxy.IContext, a int) int {
		return 0
	}, nil)

	proxy.MakeInterfaceImpl(&gen, &proxy.IContext{}, "Call1", func(ctx *proxy.IContext, a string) string {
		t.Log("called Call1")
		return "not ok"
	}, nil)

	// 调用接口方法
	gen.Call(2)
	gen.Call1("ok")

}


// TestNilImpl 测试空实现结构体方法列表
func TestNilImpl(t *testing.T) {
	gen := (*I)(nil)
	typ := reflect.TypeOf(gen).Elem()
	for i := 0; i < typ.NumMethod(); i++ {
		fmt.Println(typ.Method(i).Name, typ.Method(i).Type)
	}
}

// TestGenImpl 测试生成接口实现
func TestGenImpl(t *testing.T) {
	gen := (I)(nil)
	typ := reflect.TypeOf(&gen).Elem()
	for i := 0; i < typ.NumMethod(); i++ {
		fmt.Println(typ.Method(i).Name, typ.Method(i).Type)
	}

	genInterfaceImpl(&gen, func(data *Impl2, a int) int {
		fmt.Println("proxy")
		return 3
	})

	// 调用接口方法
	r := (gen).Call(1)

	fmt.Println("ok", r)
}

func genInterfaceImpl(i interface{}, proxy interface{}) {
	gen := hack.UnpackEFace(i).Data
	// mock接口方法
	mockfunc := reflect.ValueOf(proxy)
	ifc := *(*uintptr)(unsafe.Pointer(gen))
	fmt.Println(ifc)

	// 伪装iface
	*(*hack.Iface)(unsafe.Pointer(gen)) = hack.Iface{
		Tab: &hack.Itab{
			Fun: [99]uintptr{uintptr(mockfunc.Pointer()), uintptr(0), uintptr(0)},
		},
		Data: unsafe.Pointer(&Impl2{
			field1: "ok",
		}),
	}
	ifc = *(*uintptr)(unsafe.Pointer(gen))
	fmt.Println(ifc)
}

// TestAutoGenImpl 测试生成任意接口实现
func TestAutoGenImpl(t *testing.T) {
	logger.LogLevel = logger.DebugLevel
	logger.Log2Console(true)

	gen := (I)(nil)

	dynamicGenImpl(t, &gen)

	// 调用接口方法
	gen.Call(1)

	fmt.Println("ok")
}

func dynamicGenImpl(t *testing.T, i interface{}) {
	typ := reflect.TypeOf(i).Elem()
	for i := 0; i < typ.NumMethod(); i++ {
		fmt.Println(typ.Method(i).Name, typ.Method(i).Type)
	}

	gen := hack.UnpackEFace(i).Data

	// mock接口方法
	methodTyp := reflect.TypeOf(func(data *Impl2, a int) int {
		fmt.Println("proxy")
		return 3
	})

	mockfunc := reflect.MakeFunc(methodTyp, func(args []reflect.Value) (results []reflect.Value) {
		fmt.Println("called", args[1].Interface())
		debug.PrintStack()
		t.Log("ok")
		return []reflect.Value{reflect.ValueOf(3)}
	})
	ifc := *(*uintptr)(unsafe.Pointer(gen))
	fmt.Println(ifc)

	callStub := reflect.ValueOf(stub.MakeFuncStub).Pointer()

	mockFuncPtr := (*hack.Value)(unsafe.Pointer(&mockfunc)).Ptr
	genStub, err := stub.GenStubWithCtx(mockFuncPtr, callStub)
	if err != nil {
		panic(err)
	}

	fmt.Printf("genstub: 0x%x callstub: 0x%x\n", genStub, callStub)

	// 伪装iface
	*(*hack.Iface)(unsafe.Pointer(gen)) = hack.Iface{
		Tab: &hack.Itab{
			Fun: [99]uintptr{genStub, uintptr(0), uintptr(0)},
		},
		Data: (*hack.Value)(unsafe.Pointer(&mockfunc)).Ptr,
	}
	ifc = *(*uintptr)(unsafe.Pointer(gen))

	fmt.Println(ifc)
	fmt.Println(uintptr(getPtr(reflect.ValueOf(mockfunc.Interface()))))
	fmt.Println(mockfunc.Pointer())
}


// getPtr 获取函数的调用地址(和函数的指令地址不一样)
func getPtr(v reflect.Value) unsafe.Pointer {
	return (*hack.Value)(unsafe.Pointer(&v)).Ptr
}

func getImpl(n int) I {
	if n == 1 {
		return &Impl1{}
	} else if n == 2 {
		return &Impl2{}
	}
	return nil
}

// I 接口测试
type I interface {
	Call(int) int
	Call1(string) string
}

type Impl1 struct {
	field1 string
}

func (i Impl1) Call(a int) int {
	fmt.Println("Impl1 called ")
	return 1 + a
}

func (i Impl1) Call1(string) string {
	return "ok"
}

type Impl2 struct {
	field1 string
}

func (i Impl2) Call(a int) int {
	fmt.Println("Impl2 called ")
	return 2 + a
}

func (i Impl2) Call1(string) string {
	return "!ok"
}

// TestTraceBack 测试生成任意接口实现的traceback
func TestTraceBack(t *testing.T) {

	gen := (I)(nil)

	dynamicGenImpl(t, &gen)

	// 调用接口方法
	for i := 0; i < 1000; i++ {
		(gen).Call(1)
	}

	fmt.Println("ok")
}
