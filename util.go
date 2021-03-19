// Package mocker定义了mock的外层用户使用API定义,
// 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件汇聚了公共的工具类，比如类型转换。
package mocker

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/proxy"

	"git.code.oa.com/goom/mocker/internal/hack"
)

// currentPackageIndex 获取当前包的堆栈层次
const currentPackageIndex = 3

// currentPackageIndex 获取当前包的堆栈层次
const defaultCurrentPackageIndex = 2

// currentPackage 获取当前调用的包路径
func CurrentPackage() string {
	return currentPkg(defaultCurrentPackageIndex)
}

// currentPackage 获取当前调用的包路径
func currentPackage() string {
	return currentPkg(currentPackageIndex)
}

// currentPkg 获取调用者的包路径
func currentPkg(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	callerName := runtime.FuncForPC(pc).Name()

	if i := strings.Index(callerName, ".("); i > -1 {
		return callerName[:i]
	}

	if i := strings.LastIndex(callerName, "/"); i > -1 {
		realIndex := strings.Index(callerName[i:len(callerName)-1], ".")

		return callerName[:realIndex+i]
	}

	return callerName
}

// functionName 获取函数名称
func functionName(fnc interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name()
}

// typeName 获取类型名称
func typeName(fnc interface{}) string {
	t := reflect.TypeOf(fnc)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}

	return t.Name()
}

// inTypes 获取类型
func inTypes(isMethod bool, funTyp reflect.Type) []reflect.Type {
	skip := 0
	if isMethod {
		skip = 1
	}

	numIn := funTyp.NumIn()
	inTypes := make([]reflect.Type, numIn-skip)

	for i := 0; i < numIn-skip; i++ {
		inTypes[i] = funTyp.In(i + skip)
	}

	return inTypes
}

// outTypes 获取类型
func outTypes(funTyp reflect.Type) []reflect.Type {
	numOut := funTyp.NumOut()
	outTypes := make([]reflect.Type, numOut)

	for i := 0; i < numOut; i++ {
		outTypes[i] = funTyp.Out(i)
	}

	return outTypes
}

// I2V []interface convert to []reflect.Value
func I2V(args []interface{}, typs []reflect.Type) []reflect.Value {
	if len(args) != len(typs) {
		panic(fmt.Sprintf("args lenth mismatch,must:%d, actual:%d", len(typs), len(args)))
	}

	values := make([]reflect.Value, len(args))
	for i, a := range args {
		values[i] = toValue(a, typs[i])
	}

	return values
}

// toValue 转化为数值
func toValue(r interface{}, out reflect.Type) reflect.Value {
	v := reflect.ValueOf(r)
	if r != nil && v.Type() != out && (out.Kind() == reflect.Struct || out.Kind() == reflect.Ptr) {
		if v.Type().Size() != out.Size() {
			panic(fmt.Sprintf("type mismatch,must:%s, actual:%v", v.Type(), out))
		}
		// 类型强制转换,适用于结构体fake场景
		v = cast(v, out)
	}

	if r == nil && (out.Kind() == reflect.Interface || out.Kind() == reflect.Ptr || out.Kind() == reflect.Slice ||
		out.Kind() == reflect.Map || out.Kind() == reflect.Array || out.Kind() == reflect.Chan) {
		v = reflect.Zero(reflect.SliceOf(out).Elem())
	} else if v.Type().Kind() == reflect.Ptr &&
		v.Type() == reflect.TypeOf(&proxy.IContext{}) {
		panic("goom not support Return() API when returns mocked interface type, use Apply() API instead.")
	} else if r != nil && out.Kind() == reflect.Interface {
		ptr := reflect.New(out)

		ptr.Elem().Set(v)
		v = ptr.Elem()
	}

	return v
}

// cast 类型强制转换
func cast(v reflect.Value, typ reflect.Type) reflect.Value {
	originV := (*hack.Value)(unsafe.Pointer(&v))
	newV := reflect.NewAt(typ, originV.Ptr).Elem()
	newV1 := (*hack.Value)(unsafe.Pointer(&newV))
	v = *(*reflect.Value)(unsafe.Pointer(&hack.Value{
		Typ:  newV1.Typ,
		Ptr:  originV.Ptr,
		Flag: originV.Flag,
	}))

	return v
}

// V2I []reflect.Value convert to []interface
func V2I(args []reflect.Value, typs []reflect.Type) []interface{} {
	values := make([]interface{}, len(args))

	for i, a := range args {
		if (typs[i].Kind() == reflect.Interface || typs[i].Kind() == reflect.Ptr) && a.IsZero() {
			values[i] = nil
		} else {
			values[i] = a.Interface()
		}
	}

	return values
}
