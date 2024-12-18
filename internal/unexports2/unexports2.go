package unexports2

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/tencent/goom/erro"
	"github.com/tencent/goom/internal/hack"
)

var alignment uintptr

func init() {
	fn, err := getFunctionSymbolByName("github.com/tencent/goom/internal/unexports2.FindFuncByName")
	if err != nil {
		return
	}
	symTabAddress := uintptr(fn.Entry)
	memAddress := reflect.ValueOf(FindFuncByName).Pointer()
	alignment = memAddress - symTabAddress
}

// FindFuncByName read the symbol table at runtime
func FindFuncByName(name string) (uintptr, error) {
	fn, err := getFunctionSymbolByName(name)
	if err == nil {
		return uintptr(fn.Entry) + alignment, nil
	}
	if erro.CauseBy(err, erro.LdFlags) {
		panic(err)
	}
	return 0, err
}

// CreateFuncForCodePtr is given a code pointer and creates a function value
// that uses that pointer. The outFun argument should be a pointer to a function
// of the proper type (e.g. the address of a local variable), and will be set to
// the result function value.
func CreateFuncForCodePtr(outFuncPtr interface{}, codePtr uintptr) (*hack.Func, error) {
	outFunc := reflect.ValueOf(outFuncPtr)
	if outFunc.Kind() != reflect.Ptr {
		return nil, errors.New("func param must be ptr")
	}

	outFuncVal := outFunc.Elem()
	// Use reflect.MakeGlobalFunc to create a well-formed function value that's
	// guaranteed to be of the right type and guaranteed to be on the heap
	// (so that we can modify it). We give a nil delegate function because
	// it will never actually be called.
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil)
	// Use reflection on the reflect.Value (yep!) to grab the underling
	// function value pointer. Trying to call newFuncVal.Pointer() wouldn't
	// work because it gives the code pointer rather than the function value
	// pointer. The function value is a struct that starts with its code
	// pointer, so we can swap out the code pointer with our desired value.
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer()
	// nolint hack 用法
	funcPtr := (*hack.Func)(unsafe.Pointer(funcValuePtr))
	funcPtr.CodePtr = codePtr

	outFuncVal.Set(newFuncVal)
	return funcPtr, nil
}

// NewFuncWithCodePtr 根据类型和函数地址进行构造 reflect.Value
func NewFuncWithCodePtr(typ reflect.Type, codePtr uintptr) reflect.Value {
	pointer := unsafe.Pointer(&codePtr)
	funcVal := reflect.NewAt(typ, pointer).Elem()
	(*hack.Value)(unsafe.Pointer(&funcVal)).Flag = uintptr(reflect.Func)
	return funcVal
}
