package mocker_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"git.code.oa.com/goom/mocker"
)

// TestUnitBuilderTestSuite 测试入口
func TestUnitBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}

// BuilderTestSuite Builder测试套件
type BuilderTestSuite struct {
	suite.Suite
}

// TestUnitFunc 测试私有方法mock
func (s *BuilderTestSuite) TestUnitFunc() {
	s.Run("success", func() {
		mb := mocker.Create("")

		mb.Func("fun1").Apply(func(i int) int {
			return i * 3
		})

		mb.Func("fun2").Apply(func(i int) int {
			return i * 3
		})

		s.Equal(3, fun1(1), "fun1 mock check")
		s.Equal(3, fun2(1), "fun2 mock check")

		mb.Reset()

		s.Equal(1, fun1(1), "fun1 mock reset check")
		s.Equal(2, fun2(1), "fun1 mock reset check")
	})
}

// TestUnitMethod 测试结构体的私有方法mock
func (s *BuilderTestSuite) TestUnitMethod() {
	s.Run("success", func() {
		mb := mocker.Create("")
		mb.Struct("fake").Method("call").Apply(func(_ *fake, i int) int {
			return i * 2
		})

		f := &fake{}

		s.Equal(2, f.call(1), "call mock check")

		mb.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
	})
}

// TestUnitFuncDef 测试函数定义的mock
func (s *BuilderTestSuite) TestUnitFuncDef() {
	s.Run("success", func() {
		mb := mocker.Create("")
		mb.FuncDec((&fake{}).call).Apply(func(_ *fake, i int) int {
			return i * 2
		})

		f := &fake{}

		s.Equal(2, f.call(1), "call mock check")

		mb.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
	})
}

// TestUnitFuncReturn 测试私有方法mock return
func (s *BuilderTestSuite) TestUnitFuncReturn() {
	s.Run("success", func() {
		mb := mocker.Create("")
		mb.FuncDec(fun1).Return(3)

		s.Equal(3, fun1(1), "fun1 mock check")

		mb.Reset()

		s.Equal(1, fun1(1), "fun1 mock reset check")
	})
}

//go:noinline
func fun1(i int) int {
	return i * 1
}

//go:noinline
func fun2(i int) int {
	return i * 2
}

type fake struct{}

//go:noinline
func (f *fake) call(i int) int {
	return i
}
