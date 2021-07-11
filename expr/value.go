package expr

import (
	"fmt"
	"reflect"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/hack"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// I2V []interface convert to []reflect.Value
func I2V(args []interface{}, types []reflect.Type) []reflect.Value {
	if len(args) != len(types) {
		panic(fmt.Sprintf("args lenth mismatch,must:%d, actual:%d", len(types), len(args)))
	}

	values := make([]reflect.Value, len(args))
	for i, a := range args {
		values[i] = toValue(a, types[i])
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
func V2I(args []reflect.Value, types []reflect.Type) []interface{} {
	values := make([]interface{}, len(args))

	for i, a := range args {
		if (types[i].Kind() == reflect.Interface || types[i].Kind() == reflect.Ptr) && a.IsZero() {
			values[i] = nil
		} else {
			values[i] = a.Interface()
		}
	}

	return values
}

// ToExpr 将参数转换成[]Expr
func ToExpr(args []interface{}, types []reflect.Type) []Expr {
	if len(args) != len(types) {
		panic(fmt.Sprintf("args lenth mismatch,must:%d, actual:%d", len(types), len(args)))
	}

	exprs := make([]Expr, len(args))
	for i, a := range args {
		if expr, ok := a.(Expr); ok {
			exprs[i] = expr
		} else {
			// 默认使用equals表达式
			exprs[i] = Equals(a)
			exprs[i].Resole([]reflect.Type{types[i]})
		}
	}
	return exprs
}
