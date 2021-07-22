// Package mocker 定义了mock的外层用户使用API定义,
// 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件实现了按照参数条件进行匹配, 返回对应的mock return值,
// 支持了mocker.When(XXX).Return(YYY)的高效匹配。
package mocker

import (
	"reflect"

	"git.code.oa.com/goom/mocker/arg"
	"git.code.oa.com/goom/mocker/erro"
)

// Matcher 参数匹配接口
type Matcher interface {
	// Match 匹配执行方法
	Match(args []reflect.Value) bool
	// Result 匹配成功返回的结果
	Result() []reflect.Value
	// AddResult 添加返回结果
	AddResult([]interface{})
}

// When Mock条件匹配。
// 当参数等于指定的值时,会return对应的指定值
type When struct {
	ExportedMocker

	funcTyp        reflect.Type
	funcDef        interface{}
	isMethod       bool
	matches        []Matcher
	defaultReturns Matcher
	// curMatch 当前指定的参数匹配
	curMatch Matcher
}

// CreateWhen 构造条件判断
// args 参数条件
// defaultReturns 默认返回值
// isMethod 是否为方法类型
func CreateWhen(m ExportedMocker, funcDef interface{}, args []interface{},
	defaultReturns []interface{}, isMethod bool) (*When, error) {
	impTyp := reflect.TypeOf(funcDef)

	err := checkParams(funcDef, impTyp, args, defaultReturns, isMethod)
	if err != nil {
		return nil, err
	}

	var (
		curMatch     Matcher
		defaultMatch Matcher
	)

	if defaultReturns != nil {
		curMatch = newAlwaysMatch(defaultReturns, impTyp)
	} else if len(outTypes(impTyp)) == 0 {
		curMatch = newEmptyMatch()
	}

	defaultMatch = curMatch

	if args != nil {
		curMatch = newDefaultMatch(args, nil, isMethod, impTyp)
	}

	return &When{
		ExportedMocker: m,
		defaultReturns: defaultMatch,
		funcTyp:        impTyp,
		funcDef:        funcDef,
		isMethod:       isMethod,
		matches:        make([]Matcher, 0),
		curMatch:       curMatch,
	}, nil
}

// checkParams 检查参数
func checkParams(funcDef interface{}, impTyp reflect.Type,
	args []interface{}, returns []interface{}, isMethod bool) error {
	if returns != nil && len(returns) < impTyp.NumOut() {
		return erro.NewReturnsNotMatchError(funcDef, len(returns), impTyp.NumOut())
	}

	if isMethod {
		if args != nil && len(args)+1 < impTyp.NumIn() {
			return erro.NewArgsNotMatchError(funcDef, len(args), impTyp.NumIn()-1)
		}
	} else {
		if args != nil && len(args) < impTyp.NumIn() {
			return erro.NewArgsNotMatchError(funcDef, len(args), impTyp.NumIn())
		}
	}

	return nil
}

// NewWhen 创建默认When
func NewWhen(funTyp reflect.Type) *When {
	return &When{
		ExportedMocker: nil,
		funcTyp:        funTyp,
		matches:        make([]Matcher, 0),
		defaultReturns: nil,
		curMatch:       nil,
	}
}

// When 当参数符合一定的条件, 使用DefaultMatcher
// 入参个数必须和函数或方法参数个数一致,
// 比如: When(
//		In(3, 4), // 第一个参数是In
//		Any()) // 第二个参数是Any
func (w *When) When(args ...interface{}) *When {
	w.curMatch = newDefaultMatch(args, nil, w.isMethod, w.funcTyp)
	return w
}

// In 当参数包含其中之一, 使用ContainsMatcher
// 当参数为多个时, In的每个条件各使用一个数组表示:
// .In([]interface{}{3, Any()}, []interface{}{4, Any()})
func (w *When) In(slices ...interface{}) *When {
	w.curMatch = newContainsMatch(slices, nil, w.isMethod, w.funcTyp)
	return w
}

// Return 指定返回值
func (w *When) Return(results ...interface{}) *When {
	if w.curMatch != nil {
		w.curMatch.AddResult(results)
		w.matches = append(w.matches, w.curMatch)

		return w
	}

	if w.defaultReturns == nil {
		w.defaultReturns = newAlwaysMatch(results, w.funcTyp)
	} else {
		w.defaultReturns.AddResult(results)
	}
	return w
}

// AndReturn 指定第二次调用返回值,之后的调用以最后一个指定的值返回
func (w *When) AndReturn(results ...interface{}) *When {
	if w.curMatch == nil {
		return w.Return(results...)
	}

	w.curMatch.AddResult(results)

	return w
}

// Matches 多个条件匹配
func (w *When) Matches(matches ...arg.Pair) *When {
	if len(matches) == 0 {
		return w
	}

	for _, v := range matches {
		args, ok := v.Args.([]interface{})
		if !ok {
			args = []interface{}{v.Args}
		}

		results, ok := v.Return.([]interface{})
		if !ok {
			results = []interface{}{v.Return}
		}

		w.Return(results...)
		matcher := newDefaultMatch(args, results, w.isMethod, w.funcTyp)
		w.matches = append(w.matches, matcher)
	}

	return w
}

// Returns 按顺序依次返回值
func (w *When) Returns(rets ...interface{}) *When {
	if len(rets) == 0 {
		return w
	}

	for i, v := range rets {
		ret, ok := v.([]interface{})
		if !ok {
			ret = []interface{}{v}
		}

		if i == 0 {
			w.Return(ret...)
		} else {
			w.AndReturn(ret...)
		}
	}

	return w
}

// invoke 执行When参数匹配并返回值
func (w *When) invoke(args1 []reflect.Value) (results []reflect.Value) {
	if len(w.matches) != 0 {
		for _, c := range w.matches {
			if c.Match(args1) {
				return c.Result()
			}
		}
	}

	return w.returnDefaults()
}

// Eval 执行when子句
func (w *When) Eval(args ...interface{}) []interface{} {
	argVs := arg.I2V(args, inTypes(w.isMethod, w.funcTyp))
	resultVs := w.invoke(argVs)

	return arg.V2I(resultVs, outTypes(w.funcTyp))
}

// returnDefaults 返回默认值
func (w *When) returnDefaults() []reflect.Value {
	if w.defaultReturns == nil && w.funcTyp.NumOut() != 0 {
		panic("there is no suitable condition matched, or set default return with: mocker.Return(...)")
	}

	return w.defaultReturns.Result()
}
