package unexports2

import (
	"reflect"
	"unsafe"
)

func getFunctionAddress(function interface{}) (address uintptr, err error) {
	return reflect.ValueOf(function).Pointer(), nil
}

// NewFuncWithCodePtr 根据类型和函数地址进行构造 reflect.Value
func newFunctionWithImplementation(template interface{}, codePtr uintptr) (function interface{}, err error) {
	pointer := unsafe.Pointer(&codePtr)
	funcVal := reflect.NewAt(reflect.TypeOf(template), pointer).Elem()
	(*Value)(unsafe.Pointer(&funcVal)).Flag = uintptr(reflect.Func)
	return funcVal.Interface(), nil
}

// Value reflect.Value
type Value struct {
	Typ  *uintptr
	Ptr  unsafe.Pointer
	Flag uintptr
}
