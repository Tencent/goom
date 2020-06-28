package mocker

import (
	"reflect"
	"strconv"

	"git.code.oa.com/goom/mocker/errortype"
)

// When Mock条件匹配
type When struct {
	ExportedMocker

	funTyp         reflect.Type
	returns        []Return
	defaultReturns []interface{}
	// curArgs 当前指定的参数
	curArgs []interface{}
}

// CreateWhen 构造条件判断
// args 参数条件
// defaultReturns 默认返回值
func CreateWhen(m ExportedMocker, funcDef interface{}, args []interface{},
	defaultReturns []interface{}) (*When, error) {
	impTyp := reflect.TypeOf(funcDef)

	if defaultReturns != nil && len(defaultReturns) < impTyp.NumOut() {
		return nil, errortype.NewIllegalParamError("returns:"+strconv.Itoa(len(defaultReturns)+1), "'empty'")
	}

	if args != nil && len(args) < impTyp.NumIn() {
		return nil, errortype.NewIllegalParamError("args:"+strconv.Itoa(len(args)+1), "'empty'")
	}

	return &When{
		ExportedMocker: m,
		defaultReturns: defaultReturns,
		funTyp:         impTyp,
		curArgs:        args,
	}, nil
}

// When当参数符合一定的条件
func (w *When) When(args ...interface{}) *When {
	w.curArgs = args
	return w
}

// Return 指定返回值
func (w *When) Return(args ...interface{}) *When {
	// TODO 归档到returns
	return w
}

// Return 指定第二次调用返回值,之后的调用以最后一个指定的值返回
func (w *When) AndReturn(args ...interface{}) *When {
	// TODO 归档到returns
	return w
}

// Whens 多个条件匹配
func (w *When) Whens(argsmap map[interface{}]interface{}) *When {
	return w
}

func (w *When) invoke(args1 []reflect.Value) (results []reflect.Value) {
	if len(w.returns) != 0 {
		// TODO 支持条件判断
		return results
	}

	// 使用默认参数
	for i, r := range w.defaultReturns {
		v := reflect.ValueOf(r)

		if r == nil &&
			(w.funTyp.Out(i).Kind() == reflect.Interface || w.funTyp.Out(i).Kind() == reflect.Ptr) {
			v = reflect.Zero(reflect.SliceOf(w.funTyp.Out(i)).Elem())
		} else if r != nil && w.funTyp.Out(i).Kind() == reflect.Interface {
			ptr := reflect.New(w.funTyp.Out(i))
			ptr.Elem().Set(v)
			v = ptr.Elem()
		}

		results = append(results, v)
	}

	return results
}

type Return struct {
}
