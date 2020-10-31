package mocker

import (
	"reflect"
	"sync/atomic"

	"git.code.oa.com/goom/mocker/errortype"
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

// When Mock条件匹配
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
		defaultMatch = curMatch
	}

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
		return errortype.NewReturnsNotMatchError(funcDef, len(returns), impTyp.NumOut())
	}

	if isMethod {
		if args != nil && len(args)+1 < impTyp.NumIn() {
			return errortype.NewArgsNotMatchError(funcDef, len(args), impTyp.NumIn()-1)
		}
	} else {
		if args != nil && len(args) < impTyp.NumIn() {
			return errortype.NewArgsNotMatchError(funcDef, len(args), impTyp.NumIn())
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

// When 当参数符合一定的条件
func (w *When) When(args ...interface{}) *When {
	w.curMatch = newDefaultMatch(args, nil, w.isMethod, w.funcTyp)
	return w
}

// In 当参数包含其中之一
func (w *When) In(slices ...interface{}) *When {
	w.curMatch = newContainsMatch(slices, nil, w.isMethod, w.funcTyp)
	return w
}

// Matcher 指定返回值
func (w *When) Return(results ...interface{}) *When {
	if w.curMatch == nil {
		w.defaultReturns = newAlwaysMatch(results, w.funcTyp)
		return w
	}

	w.curMatch.AddResult(results)
	w.matches = append(w.matches, w.curMatch)

	return w
}

// Matcher 指定第二次调用返回值,之后的调用以最后一个指定的值返回
func (w *When) AndReturn(results ...interface{}) *When {
	if w.curMatch == nil {
		return w.Return(results...)
	}

	w.curMatch.AddResult(results)

	return w
}

// Returns 多个条件匹配
func (w *When) Returns(resultsmap map[interface{}]interface{}) *When {
	if len(resultsmap) == 0 {
		return w
	}

	for k, v := range resultsmap {
		args, ok := k.([]interface{})
		if !ok {
			args = []interface{}{k}
		}

		results, ok := v.([]interface{})
		if !ok {
			results = []interface{}{v}
		}

		matcher := newDefaultMatch(args, results, w.isMethod, w.funcTyp)
		w.matches = append(w.matches, matcher)
	}

	return w
}

// invoke invoke
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
	argVs := I2V(args, inTypes(w.isMethod, w.funcTyp))
	resultVs := w.invoke(argVs)

	return V2I(resultVs, outTypes(w.funcTyp))
}

// returnDefaults 返回默认值
func (w *When) returnDefaults() []reflect.Value {
	if w.defaultReturns == nil && w.funcTyp.NumOut() != 0 {
		panic("default returns not set.")
	}

	return w.defaultReturns.Result()
}

// BaseMatcher 参数匹配基类
type BaseMatcher struct {
	results [][]reflect.Value
	curNum  int32
	funTyp  reflect.Type
}

//newBaseMatcher 创建新参数匹配基类
func newBaseMatcher(results []interface{}, funTyp reflect.Type) *BaseMatcher {
	resultVs := make([][]reflect.Value, 0)
	if results != nil {
		// TODO results check
		resultVs = append(resultVs, I2V(results, outTypes(funTyp)))
	}

	return &BaseMatcher{
		results: resultVs,
		curNum:  0,
		funTyp:  funTyp,
	}
}

//noLint
func (c *BaseMatcher) Result() []reflect.Value {
	if len(c.results) <= 1 {
		return c.results[c.curNum]
	}

	curNum := atomic.LoadInt32(&c.curNum)
	if len := len(c.results); curNum >= int32(len) {
		return c.results[len-1]
	}

	atomic.AddInt32(&c.curNum, 1)

	return c.results[curNum]
}

// AddResult 添加结果
func (c *BaseMatcher) AddResult(results []interface{}) {
	// TODO results check
	c.results = append(c.results, I2V(results, outTypes(c.funTyp)))
}

// DefaultMatcher 参数匹配
type DefaultMatcher struct {
	*BaseMatcher

	isMethod bool
	args     []reflect.Value
}

//newDefaultMatch 创建新参数匹配
func newDefaultMatch(args []interface{}, results []interface{}, isMethod bool, funTyp reflect.Type) *DefaultMatcher {
	argVs := I2V(args, inTypes(isMethod, funTyp))

	return &DefaultMatcher{
		args:        argVs,
		BaseMatcher: newBaseMatcher(results, funTyp),
		isMethod:    isMethod,
	}
}

//Match 判断是否匹配
func (c *DefaultMatcher) Match(args []reflect.Value) bool {
	if c.isMethod {
		if len(args) != len(c.args)+1 {
			return false
		}
	} else {
		if len(args) != len(c.args) {
			return false
		}
	}

	skip := 0
	if c.isMethod {
		skip = 1
	}

	for i, arg := range c.args {
		if !equal(arg, args[i+skip]) {
			return false
		}
	}

	return true
}

// ContainsMatcher 包含类型的参数匹配
type ContainsMatcher struct {
	*BaseMatcher

	args     [][]reflect.Value
	isMethod bool
}

//newContainsMatch 创建新的包含类型的参数匹配
func newContainsMatch(args []interface{}, results []interface{}, isMethod bool, funTyp reflect.Type) *ContainsMatcher {
	argVs := make([][]reflect.Value, 0)

	for _, v := range args {
		arg, ok := v.([]interface{})
		if !ok {
			arg = []interface{}{v}
		}
		// TODO results check
		values := I2V(arg, inTypes(isMethod, funTyp))
		argVs = append(argVs, values)
	}

	return &ContainsMatcher{
		args:        argVs,
		BaseMatcher: newBaseMatcher(results, funTyp),
		isMethod:    isMethod,
	}
}

//Match 判断是否匹配
func (c *ContainsMatcher) Match(args []reflect.Value) bool {
outer:
	for _, one := range c.args {
		if c.isMethod {
			if len(args) != len(one)+1 {
				continue
			}
		} else {
			if len(args) != len(one) {
				continue
			}
		}

		skip := 0
		if c.isMethod {
			skip = 1
		}
		for i, arg := range one {
			if !equal(arg, args[i+skip]) {
				continue outer
			}
		}

		return true
	}

	return false
}

// AlwaysMatcher 默认匹配
type AlwaysMatcher struct {
	*BaseMatcher
}

// newAlwaysMatch 创建新的默认匹配
func newAlwaysMatch(results []interface{}, funTyp reflect.Type) *AlwaysMatcher {
	if results == nil {
		return nil
	}

	return &AlwaysMatcher{
		BaseMatcher: newBaseMatcher(results, funTyp),
	}
}

// Match 总是匹配
func (c *AlwaysMatcher) Match(args []reflect.Value) bool {
	return true
}
