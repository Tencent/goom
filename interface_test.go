package mocker_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"git.code.oa.com/goom/mocker"
)

// TestInterfaceCall 测试接口调用
func TestInterfaceCall(t *testing.T) {
	i := getImpl(1)
	i.Call()
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

	// mock接口方法
	mockfunc := reflect.ValueOf(func(data *Impl2) int {
		fmt.Println("proxy")
		return 3
	})

	ifc := *(*uintptr)(unsafe.Pointer(&gen))
	fmt.Println(ifc)

	// 伪装iface
	*(*iface)(unsafe.Pointer(&gen)) = iface{
		tab: &itab{
			fun: [3]uintptr{uintptr(mockfunc.Pointer()), uintptr(0), uintptr(0)},
		},
		data: unsafe.Pointer(&Impl2{
			field1: "ok",
		}),
	}

	ifc = *(*uintptr)(unsafe.Pointer(&gen))
	fmt.Println(ifc)

	// 调用接口方法
	r := (gen).Call()

	fmt.Println("ok", r)
}

// TestAutoGenImpl 测试生成任意接口实现
func TestAutoGenImpl(t *testing.T) {
	gen := (I)(nil)
	typ := reflect.TypeOf(&gen).Elem()
	for i := 0; i < typ.NumMethod(); i++ {
		fmt.Println(typ.Method(i).Name, typ.Method(i).Type)
	}

	// mock接口方法
	methodTyp := reflect.TypeOf(func(data *Impl2) int {
		fmt.Println("proxy")
		return 3
	})

	mockfunc := reflect.MakeFunc(methodTyp, func(args []reflect.Value) (results []reflect.Value) {
		fmt.Println("called")
		return []reflect.Value{reflect.ValueOf(3)}
	})

	ifc := *(*uintptr)(unsafe.Pointer(&gen))
	fmt.Println(ifc)

	// 伪装iface
	*(*iface)(unsafe.Pointer(&gen)) = iface{
		tab: &itab{
			fun: [3]uintptr{uintptr(reflect.ValueOf(mocker.InterfaceCallStub).Pointer()), uintptr(0), uintptr(0)},
		},
		data: (*value)(unsafe.Pointer(&mockfunc)).ptr,
	}

	ifc = *(*uintptr)(unsafe.Pointer(&gen))
	fmt.Println(ifc)

	fmt.Println(uintptr(getPtr(reflect.ValueOf(mockfunc.Interface()))))
	fmt.Println(mockfunc.Pointer())

	// 调用接口方法
	(gen).Call()

	fmt.Println("ok")
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type itab struct {
	inter *uintptr
	_type *uintptr
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [3]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type value struct {
	_   uintptr
	ptr unsafe.Pointer
}

type makeFuncImpl struct {
	code   uintptr
	stack  *uintptr // ptrmap for both args and results
	argLen uintptr  // just args
	ftyp   *uintptr
}

// getPtr 获取函数的调用地址(和函数的指令地址不一样)
func getPtr(v reflect.Value) unsafe.Pointer {
	return (*value)(unsafe.Pointer(&v)).ptr
}

type enhanceFuncType struct {
}

//go:noinline
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
	Call() int
}

type Impl1 struct {
	field1 string
}

func (i Impl1) Call() int {
	fmt.Println("Impl1 called ")
	return 1
}

type Impl2 struct {
	field1 string
}

func (i Impl2) Call() int {
	fmt.Println("Impl2 called ")
	return 2
}
