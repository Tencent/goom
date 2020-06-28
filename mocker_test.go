package mocker_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"git.code.oa.com/goom/mocker"
)

// TestUnitBuilderTestSuite 测试入口
func TestUnitBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(MockerTestSuite))
}

// MockerTestSuite Builder测试套件
type MockerTestSuite struct {
	suite.Suite
}

// TestUnitFunc 测试函数mock return
func (s *MockerTestSuite) TestUnitFunc() {
	s.Run("success", func() {
		mb := mocker.Create("")
		mb.Func(fun1).Return(3)

		s.Equal(3, fun1(1), "fun1 mock check")

		mb.Reset()

		s.Equal(1, fun1(1), "fun1 mock reset check")
	})
}

// TestUnitMethod 测试结构体的方法mock return
func (s *MockerTestSuite) TestUnitMethod() {
	s.Run("success", func() {
		mb := mocker.Create("")
		mb.Struct(&fake{}).Method("Call").Return(5)

		f := &fake{}

		s.Equal(5, f.Call(1), "call mock check")

		mb.Reset()

		s.Equal(1, f.Call(1), "call mock reset check")
	})
}

// TestUnitUnexportMethod 测试结构体的未导出方法mock apply
func (s *MockerTestSuite) TestUnitUnexportMethod() {
	s.Run("success", func() {
		mb := mocker.Create("")
		mb.Struct(&fake{}).UnexportedMethod("call").Apply(func(_ *fake, i int) int {
			return i * 2
		})

		f := &fake{}

		s.Equal(2, f.call(1), "call mock check")

		mb.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
	})
}

// TestUnitFunc 测试未导出函数mock
func (s *MockerTestSuite) TestUnitUnexportedFunc() {
	s.Run("success", func() {
		mb := mocker.Create("git.code.oa.com/goom/mocker_test")

		mb.UnexportedFunc("fun1").Apply(func(i int) int {
			return i * 3
		})

		mb.UnexportedFunc("fun2").Apply(func(i int) int {
			return i * 3
		})

		s.Equal(3, fun1(1), "fun1 mock check")
		s.Equal(3, fun2(1), "fun2 mock check")

		mb.Reset()

		s.Equal(1, fun1(1), "fun1 mock reset check")
		s.Equal(2, fun2(1), "fun1 mock reset check")
	})
}

// TestUnitUnexportMethod 测试未导出结构体的方法mock apply
func (s *MockerTestSuite) TestUnitUnexportStruct() {
	s.Run("success", func() {
		// 指定包名
		mb := mocker.Create("git.code.oa.com/goom/mocker_test")
		mb.UnexportedStruct("fake").Method("call").Apply(func(_ *fake, i int) int {
			return i * 2
		})

		f := &fake{}

		s.Equal(2, f.call(1), "call mock check")

		mb.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
	})
}

// TestUnitAny 测试任意函数定义的mock
func (s *MockerTestSuite) TestUnitAny() {
	s.Run("success", func() {
		mb := mocker.Create("")
		// 指定: &fake{}).call,此方式不支持return
		mb.Func((&fake{}).call).Apply(func(_ *fake, i int) int {
			return i * 2
		})

		f := &fake{}

		s.Equal(2, f.call(1), "call mock check")

		mb.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
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
func (f *fake) Call(i int) int {
	return i
}

//go:noinline
func (f *fake) call(i int) int {
	return i
}
