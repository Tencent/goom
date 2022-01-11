package arg

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"unsafe"

	"git.code.oa.com/goom/mocker/internal/hack"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// I2V []interface convert to []reflect.Value
func I2V(args []interface{}, types []reflect.Type) []reflect.Value {
	if len(args) != len(types) {
		panic(fmt.Sprintf("args lenth mismatch,must: %d, actual: %d", len(types), len(args)))
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
			panic(fmt.Sprintf("type mismatch,must: %s, actual: %v", v.Type(), out))
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
	newVHack := (*hack.Value)(unsafe.Pointer(&newV))
	v = *(*reflect.Value)(unsafe.Pointer(&hack.Value{
		Typ:  newVHack.Typ,
		Ptr:  originV.Ptr,
		Flag: originV.Flag,
	}))

	return v
}

// V2I []reflect.Value convert to []interface
func V2I(args []reflect.Value, types []reflect.Type) []interface{} {
	values := make([]interface{}, len(args))

	for i, a := range args {
		if (types[i].Kind() == reflect.Interface || types[i].Kind() == reflect.Ptr) && isZero(a) {
			values[i] = nil
		} else {
			values[i] = a.Interface()
		}
	}

	return values
}

// SprintV []reflect.Value print to string
func SprintV(args []reflect.Value) string {
	s := make([]string, 0, len(args))

	for _, a := range args {
		if (a.Kind() == reflect.Interface || a.Kind() == reflect.Ptr) && isZero(a) {
			s = append(s, "nil")
		} else {
			s = append(s, fmt.Sprintf("%v", a.Interface()))
		}
	}

	return strings.Join(s, ",")
}

// ToExpr 将参数转换成[]Expr
func ToExpr(args []interface{}, types []reflect.Type) ([]Expr, error) {
	if len(args) != len(types) {
		return nil, fmt.Errorf("args lenth mismatch,must: %d, actual: %d", len(types), len(args))
	}

	// TODO results check
	exprs := make([]Expr, len(args))
	for i, a := range args {
		if expr, ok := a.(Expr); ok {
			exprs[i] = expr
		} else {
			// 默认使用equals表达式
			exprs[i] = Equals(a)
		}
		err := exprs[i].Resolve([]reflect.Type{types[i]})
		if err != nil {
			return nil, err
		}
	}
	return exprs, nil
}

// isZero reports whether v is the zero value for its type.
// It panics if the argument is invalid.
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return true
	}
}
