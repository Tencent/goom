package mocker

import (
	"reflect"
	"strconv"
	"sync/atomic"

	"git.code.oa.com/goom/mocker/errortype"
)

// Matcher 参数匹配接口
type Matcher interface {
	Match(args []reflect.Value) bool
	Result() []reflect.Value
	AddResult([]interface{})
}

// When Mock条件匹配
type When struct {
	ExportedMocker

	funTyp         reflect.Type
	matches        []Matcher
	defaultReturns []interface{}
	// curMatch 当前指定的参数匹配
	curMatch Matcher
}

// CreateWhen 构造条件判断
// args 参数条件
// defaultReturns 默认返回值
func CreateWhen(m ExportedMocker, funcDef interface{}, args []interface{},
	defaultReturns []interface{}) (*When, error) {
	impTyp := reflect.TypeOf(funcDef)

	if defaultReturns != nil && len(defaultReturns) < impTyp.NumOut() {
		return nil, errortype.NewIllegalParamError("matches:"+
			strconv.Itoa(len(defaultReturns)+1), "'empty'")
	}

	if args != nil && len(args) < impTyp.NumIn() {
		return nil, errortype.NewIllegalParamError("args:"+
			strconv.Itoa(len(args)+1), "'empty'")
	}

	return &When{
		ExportedMocker: m,
		defaultReturns: defaultReturns,
		funTyp:         impTyp,
		matches:        make([]Matcher, 0),
		curMatch:       newDefaultMatch(args, nil),
	}, nil
}

// NewWhen 创建默认When
func NewWhen(funTyp reflect.Type) *When {
	return &When{
		ExportedMocker: nil,
		funTyp:         funTyp,
		matches:        make([]Matcher, 0),
		defaultReturns: nil,
		curMatch:       nil,
	}
}

// When 当参数符合一定的条件
func (w *When) When(args ...interface{}) *When {
	w.curMatch = newDefaultMatch(args, nil)
	return w
}

// In 当参数包含其中之一
func (w *When) In(slices ...interface{}) *When {
	w.curMatch = newContainsMatch(slices, nil)
	return w
}

// Matcher 指定返回值
func (w *When) Return(results ...interface{}) *When {
	if w.curMatch == nil {
		w.defaultReturns = results
		return w
	}

	w.curMatch.AddResult(results)
	w.matches = append(w.matches, w.curMatch)
	return w
}

// Matcher 指定第二次调用返回值,之后的调用以最后一个指定的值返回
func (w *When) AndReturn(results ...interface{}) *When {
	if w.curMatch == nil {
		return w
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

		matcher := newDefaultMatch(args, results)
		w.matches = append(w.matches, matcher)
	}

	return w
}

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
	argVs := I2V(args)
	resultVs := w.invoke(argVs)
	return V2I(resultVs)
}

func (w *When) returnDefaults() []reflect.Value {
	if w.defaultReturns == nil {
		panic("default whens not set.")
	}

	var results []reflect.Value
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

// BaseMatcher 参数匹配基类
type BaseMatcher struct {
	results [][]reflect.Value
	curNum  int32
}

func newBaseMatcher(results []interface{}) *BaseMatcher {
	resultVs := make([][]reflect.Value, 0)
	if results != nil {
		// TODO results check
		resultVs = append(resultVs, I2V(results))
	}
	return &BaseMatcher{
		results: resultVs,
		curNum:  0,
	}
}

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

func (c *BaseMatcher) AddResult(results []interface{}) {
	// TODO results check
	c.results = append(c.results, I2V(results))
}

// Matcher 参数匹配
type DefaultMatcher struct {
	*BaseMatcher

	args []reflect.Value
}

func newDefaultMatch(args []interface{}, results []interface{}) *DefaultMatcher {
	argVs := I2V(args)
	return &DefaultMatcher{
		args:        argVs,
		BaseMatcher: newBaseMatcher(results),
	}
}

func (c *DefaultMatcher) Match(args []reflect.Value) bool {
	if len(args) != len(c.args) {
		return false
	}

	for i, arg := range c.args {
		if !equal(arg, args[i]) {
			return false
		}
	}

	return true
}

// ContainsMatcher 包含类型的参数匹配
type ContainsMatcher struct {
	*BaseMatcher

	args [][]reflect.Value
}

func newContainsMatch(args []interface{}, results []interface{}) *ContainsMatcher {
	argVs := make([][]reflect.Value, 0)

	for _, v := range args {
		arg, ok := v.([]interface{})
		if !ok {
			arg = []interface{}{v}
		}
		// TODO results check
		values := I2V(arg)
		argVs = append(argVs, values)
	}

	return &ContainsMatcher{
		args:        argVs,
		BaseMatcher: newBaseMatcher(results),
	}
}

func (c *ContainsMatcher) Match(args []reflect.Value) bool {
outer:
	for _, one := range c.args {
		if len(args) != len(one) {
			continue
		}
		for i, arg := range one {
			if !equal(arg, args[i]) {
				continue outer
			}
		}
		return true
	}
	return false
}
