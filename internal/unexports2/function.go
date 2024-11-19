package unexports2

import (
	"reflect"
	"unsafe"
)

func getFunctionAddress(function interface{}) (address uintptr, err error) {
	return reflect.ValueOf(function).Pointer(), nil
	//rv := reflect.ValueOf(function)
	//if err = MakeAddressable(&rv); err != nil {
	//	return
	//}
	//pFunc := (*unsafe.Pointer)(unsafe.Pointer(rv.UnsafeAddr()))
	//address = uintptr(*pFunc)
	//return
}

func newFunctionWithImplementation1(template interface{}, implementationPtr uintptr) (function interface{}, err error) {
	rFunc := reflect.MakeFunc(reflect.TypeOf(template), nil)
	if err = MakeAddressable(&rFunc); err != nil {
		return
	}
	pFunc := (*unsafe.Pointer)(unsafe.Pointer(rFunc.UnsafeAddr()))
	*pFunc = unsafe.Pointer(uintptr(implementationPtr))
	function = rFunc.Interface()
	return
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
