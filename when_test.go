package mocker_test

import (
	"git.code.oa.com/goom/mocker"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

// TestUnitWhenTestSuite 测试入口
func TestUnitWhenTestSuite(t *testing.T) {
	suite.Run(t, new(WhenTestSuite))
}

// MockerTestSuite Builder测试套件
type WhenTestSuite struct {
	suite.Suite
}


func simple(a int) int {
	return 0
}


type Arg struct {
	field1 string
}

type Result struct {
	field1 int
}

func complex(a Arg) Result {
	return Result{0}
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
		when.Return(-1).WhenContains(1, 2).Return(5)

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


// TestNil 测试空参数
func (s *WhenTestSuite) TestComplex() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(complex))
		when.Return(Result{}).When(Arg{field1:"ok"}).Return(Result{0}).
			When(Arg{field1:"not ok"}).Return(Result{-1})

		s.Equal(Result{0}, when.Eval(Arg{field1:"ok"})[0], "when result check")
		s.Equal(Result{-1}, when.Eval(Arg{field1:"not ok"})[0], "when result check")
		s.Equal(Result{}, when.Eval(Arg{field1:"other"})[0], "when result check")
	})
}

// TestNil 测试空参数
func (s *WhenTestSuite) TestNil() {
	s.Run("success", func() {
		when := mocker.NewWhen(reflect.TypeOf(simple))
		when.Return(-1).When(1).Return(nil)

		//s.Equal(5, when.Eval(1)[0], "when result check")
	})
}
