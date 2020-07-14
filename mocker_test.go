package mocker_test

import (
	"fmt"
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

// TestUnitFuncApply 测试函数mock apply
func (s *MockerTestSuite) TestUnitFuncApply() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).Apply(func(int) int {
			return 3
		})

		s.Equal(3, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestUnitFuncReturn 测试函数mock return
func (s *MockerTestSuite) TestUnitFuncReturn() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).When(1).Return(3)

		s.Equal(3, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestUnitUnexportedFuncApply 测试未导出函数mock apply
func (s *MockerTestSuite) TestUnitUnexportedFuncApply() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Pkg("git.code.oa.com/goom/mocker_test").ExportFunc("foo").Apply(func(i int) int {
			return i * 3
		})

		s.Equal(3, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestUnitUnexportedFuncReturn 测试未导出函数mock return
func (s *MockerTestSuite) TestUnitUnexportedFuncReturn() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Pkg("git.code.oa.com/goom/mocker_test").ExportFunc("foo").As(func(i int) int {
			return i * 3
		}).Return(3)

		s.Equal(3, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestUnitMethodApply 测试结构体的方法mock apply
func (s *MockerTestSuite) TestUnitMethodApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&fake{}).Method("Call").Apply(func(*fake, int) int {
			return 5
		})

		f := &fake{}

		s.Equal(5, f.Call(1), "call mock check")

		mock.Reset()

		s.Equal(1, f.Call(1), "call mock reset check")
	})
}

// TestUnitMethodReturn 测试结构体的方法mock return
func (s *MockerTestSuite) TestUnitMethodReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&fake{}).Method("Call").Return(5)

		f := &fake{}

		s.Equal(5, f.Call(1), "call mock check")

		mock.Reset()

		s.Equal(1, f.Call(1), "call mock reset check")
	})
}

// TestUnitUnexportMethodApply 测试结构体的未导出方法mock apply
func (s *MockerTestSuite) TestUnitUnexportedMethodApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&fake{}).ExportMethod("call").Apply(func(_ *fake, i int) int {
			return i * 2
		})

		f := &fake{}

		s.Equal(2, f.call(1), "call mock check")

		mock.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
	})
}

// TestUnitUnexportedMethodReturn 测试结构体的未导出方法mock return
func (s *MockerTestSuite) TestUnitUnexportedMethodReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&fake{}).ExportMethod("call").As(func(_ *fake, i int) int {
			return i * 2
		}).Return(6)

		f := &fake{}

		s.Equal(6, f.call(1), "call mock check")

		mock.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
	})
}

// TestUnitUnexportStruct 测试未导出结构体的方法mock apply
func (s *MockerTestSuite) TestUnitUnexportStruct() {
	s.Run("success", func() {
		// 指定包名
		mock := mocker.Create()
		mock.Pkg("git.code.oa.com/goom/mocker_test").ExportStruct("*fake").
			Method("call").Apply(func(_ *fake, i int) int {
			return i * 2
		})

		f := &fake{}

		s.Equal(2, f.call(1), "call mock check")

		mock.Reset()

		s.Equal(1, f.call(1), "call mock reset check")
	})
}

// TestCallOrigin 测试调用原函数mock return
func (s *MockerTestSuite) TestCallOrigin() {
	s.Run("success", func() {
		mock := mocker.Create()

		// 定义原函数,用于占位,实际不会执行该函数体
		var origin = func(i int) int {
			fmt.Println("origin func placeholder")
			return 0 + i
		}

		mock.Func(foo).Origin(&origin).Apply(func(i int) int {
			originResult := origin(i)
			return originResult + 100
		})

		s.Equal(101, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

func (s *MockerTestSuite) TestMultiReturn() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).When(1).Return(3).AndReturn(2)

		s.Equal(3, foo(1), "foo mock check")
		s.Equal(2, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

//go:noinline
func foo(i int) int {
	return i * 1
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
