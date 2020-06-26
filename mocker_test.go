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


// TestUnitFunc 测试私有方法mock return
func (s *BuilderTestSuite) TestUnitFunc() {
	s.Run("success", func() {
		mb := mocker.Create("")
		mb.Func(fun1).Return(3)

		s.Equal(3, fun1(1), "fun1 mock check")

		mb.Reset()

		s.Equal(1, fun1(1), "fun1 mock reset check")
	})
}

// TestUnitMethod 测试结构体的私有方法mock return
func (s *BuilderTestSuite) TestUnitMethod() {
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
func (s *BuilderTestSuite) TestUnitUnexportMethod() {
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

// TestUnitFunc 测试私有方法mock
func (s *BuilderTestSuite) TestUnitUnexportedFunc() {
	s.Run("success", func() {
		mb := mocker.Create("git.code.oa.com/goom/mocker_test")

		mb.UnexportF("fun1").Apply(func(i int) int {
			return i * 3
		})

		mb.UnexportF("fun2").Apply(func(i int) int {
			return i * 3
		})

		s.Equal(3, fun1(1), "fun1 mock check")
		s.Equal(3, fun2(1), "fun2 mock check")

		mb.Reset()

		s.Equal(1, fun1(1), "fun1 mock reset check")
		s.Equal(2, fun2(1), "fun1 mock reset check")
	})
}

// TestUnitUnexportMethod 测试结构体的未导出方法mock apply
func (s *BuilderTestSuite) TestUnitUnexportStruct() {
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
func (s *BuilderTestSuite) TestUnitAny() {
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
