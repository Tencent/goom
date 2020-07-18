package proxy_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/hack"

	"git.code.oa.com/goom/mocker/internal/proxy"
)

// TestInterfaceCall 测试接口调用
func TestInterfaceCall(t *testing.T) {
	i := getImpl(1)
	i.Call(0)
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

	genInterfaceImpl(&gen)

	// 调用接口方法
	r := (gen).Call(1)

	fmt.Println("ok", r)
}

func genInterfaceImpl(i interface{}) {
	gen := hack.UnpackEFace(i).Data
	// mock接口方法
	mockfunc := reflect.ValueOf(func(data *Impl2, a int) int {
		fmt.Println("proxy")
		return 3
	})
	ifc := *(*uintptr)(unsafe.Pointer(gen))
	fmt.Println(ifc)
	// 伪装iface
	*(*hack.Iface)(unsafe.Pointer(gen)) = hack.Iface{
		Tab: &hack.Itab{
			Fun: [3]uintptr{uintptr(mockfunc.Pointer()), uintptr(0), uintptr(0)},
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
	gen := (I)(nil)

	dynamicGenImpl(&gen)

	// 调用接口方法
	(gen).Call(1)

	fmt.Println("ok")
}

func dynamicGenImpl(i interface{}) {
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
		return []reflect.Value{reflect.ValueOf(3)}
	})
	ifc := *(*uintptr)(unsafe.Pointer(gen))
	fmt.Println(ifc)

	// 伪装iface
	*(*hack.Iface)(unsafe.Pointer(gen)) = hack.Iface{
		Tab: &hack.Itab{
			Fun: [3]uintptr{uintptr(reflect.ValueOf(proxy.InterfaceCallStub).Pointer()), uintptr(0), uintptr(0)},
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
