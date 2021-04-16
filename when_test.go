// Package mocker_test 对mocker包的测试
// 当前文件实现了对when.go的单测
package mocker_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	"git.code.oa.com/goom/mocker"
)

// TestUnitWhenTestSuite 测试入口
func TestUnitWhenTestSuite(t *testing.T) {
	suite.Run(t, new(WhenTestSuite))
}

// MockerTestSuite Builder测试套件
type WhenTestSuite struct {
	suite.Suite
}

// simple 普通函数
func simple(int) int {
	return 0
}

// Arg 普通参数
type Arg struct {
	field1 string
}

// Result 普通返回结果
type Result struct {
	field1 int
}

// complex 复杂返回结果函数
func complex(Arg) Result {
	return Result{0}
}

// complex1 复杂带指针的返回结果函数
func complex1(Arg) *Result {
	return &Result{0}
}

// Struct for结构体方法When
type Struct struct{}

// Div 除法操作
//go:noinline
func (s *Struct) Div(a int, b int) int {
	return a / b
}

// StructOuter 嵌套结构外层
type StructOuter struct {
}

// Compute 中会调用sub运算
func (s *StructOuter) Compute(a int, b int) int {
	diver := new(Struct)
	res := diver.Div(a, b)

	return res
}

// TestWhen 测试简单参数匹配
func (s *WhenTestSuite) TestWhen() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(simple))
		when.When(1).Return(2)

		s.Equal(2, when.Eval(1)[0], "when result check")
	})
}

// TestWhenAndReturn 多次返回不同的值
func (s *WhenTestSuite) TestWhenAndReturn() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(simple))
		when.When(1).Return(2).AndReturn(3)

		s.Equal(2, when.Eval(1)[0], "when result check")
		s.Equal(3, when.Eval(1)[0], "when result check")
		s.Equal(3, when.Eval(1)[0], "when result check")
	})
}

// TestWhenContains 任意一个配
func (s *WhenTestSuite) TestWhenContains() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(simple))
		when.Return(-1).In(1, 2).Return(5)

		s.Equal(5, when.Eval(1)[0], "when result check")
		s.Equal(5, when.Eval(1)[0], "when result check")
		s.Equal(-1, when.Eval(0)[0], "when result check")
	})
}

// TestReturns 测试批量设置条件
func (s *WhenTestSuite) TestReturns() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(simple))
		when.Return(-1).Returns(map[interface{}]interface{}{
			1: 5,
			2: 5,
			3: 6,
		})

		s.Equal(5, when.Eval(1)[0], "when result check")
		s.Equal(5, when.Eval(2)[0], "when result check")
		s.Equal(6, when.Eval(3)[0], "when result check")
	})
}

// TestNil 测试复杂参数
func (s *WhenTestSuite) TestComplex() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(complex))
		when.Return(Result{}).When(Arg{field1: "ok"}).Return(Result{0}).
			When(Arg{field1: "not ok"}).Return(Result{-1})

		s.Equal(Result{0}, when.Eval(Arg{field1: "ok"})[0], "when result check")
		s.Equal(Result{-1}, when.Eval(Arg{field1: "not ok"})[0], "when result check")
		s.Equal(Result{}, when.Eval(Arg{field1: "other"})[0], "when result check")
	})
}

// TestNil 测试空参数
func (s *WhenTestSuite) TestNil() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(complex1))
		when.Return(nil)

		s.Equal(nil, when.Eval(Arg{})[0], "when return nil check")
	})
}

// TestMethodWhen 方法参数条件匹配
func (s *WhenTestSuite) TestMethodWhen() {
	s.Run("success", func() {
		structOuter := new(StructOuter)
		struct1 := new(Struct)
		m := mocker.Create()

		// 直接mock方法的返回值
		m.Struct(struct1).Method("Div").Return(100)
		s.Equal(100, structOuter.Compute(2, 1), "method when check")

		m.Reset()
		m.Struct(struct1).Method("Div").Return(50)
		s.Equal(50, structOuter.Compute(2, 1), "method when check")

		m.Struct(struct1).Method("Div").When(3, 4).Return(100)
		m.Struct(struct1).Method("Div").When(4, 4).Return(200)
		s.Equal(100, structOuter.Compute(3, 4), "method when check")
		s.Equal(200, structOuter.Compute(4, 4), "method when check")

		// mock方法的替换方法
		m.Struct(struct1).Method("Div").Apply(func(_ *Struct, a int, b int) int {
			return a/b + 1
		})
		s.Equal(3, structOuter.Compute(2, 1), "method when check")
	})
}
