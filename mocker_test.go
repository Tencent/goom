package mocker_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

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

// TestUnitInterfaceApply 测试接口mock apply
func (s *MockerTestSuite) TestUnitInterfaceApply() {
	s.Run("success", func() {
		mock := mocker.Create()

		i := (I)(nil)

		mock.Interface(&i).Method("Call").Apply(func(ctx *mocker.IContext, i int) int {
			return 3
		})
		mock.Interface(&i).Method("Call1").Apply(func(ctx *mocker.IContext, i string) string {
			return "ok"
		})

		s.Equal(3, i.Call(1), "interface mock check")
		s.Equal("ok", i.Call1(""), "interface mock check")

		mock.Reset()

		s.Equal(nil, i, "interface mock reset check")
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

// TestUnitInterfaceReturn 测试接口mock return
func (s *MockerTestSuite) TestUnitInterfaceReturn() {
	s.Run("success", func() {
		mock := mocker.Create()

		i := (I)(nil)

		mock.Interface(&i).Method("Call").As(func(ctx *mocker.IContext, i int) int {
			return 0
		}).When(1).Return(3)
		mock.Interface(&i).Method("Call1").As(func(ctx *mocker.IContext, s string) string {
			return ""
		}).When("").Return("ok")
		mock.Interface(&i).Method("call2").As(func(ctx *mocker.IContext, i int32) int32 {
			return 0
		}).Return(int32(5))

		s.Equal(3, i.Call(1), "interface mock check")
		s.Equal("ok", i.Call1(""), "interface mock check")
		s.Equal(int32(5), i.call2(0), "interface mock check")

		mock.Reset()

		s.Equal(nil, i, "interface mock reset check")
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
		mock.Struct(&fake{}).Method("Call").Return(5).AndReturn(6)
		mock.Struct(&fake{}).Method("Call2").Return(7).AndReturn(8)

		f := &fake{}

		s.Equal(5, f.Call(1), "call mock check")
		s.Equal(6, f.Call(1), "call mock check")
		s.Equal(7, f.Call2(1), "call mock check")
		s.Equal(8, f.Call2(1), "call mock check")

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

// TestUnitFuncTwiceApply 测试函数mock apply多次
func (s *MockerTestSuite) TestUnitFuncTwiceApply() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).When(1).Return(3)
		mock.Func(foo).When(2).Return(6)
		s.Equal(3, foo(1), "foo mock check")
		s.Equal(6, foo(2), "foo mock check")
		mock.Reset()

		mock.Func(foo).When(1).Return(2)
		s.Equal(2, foo(1), "foo mock reset check")
		mock.Reset()
	})
}

// TestUnitDefaultReturn 测试函数mock返回默认值
func (s *MockerTestSuite) TestUnitDefaultReturn() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).Return(3).AndReturn(4)
		mock.Func(foo).Return(5).AndReturn(6)
		s.Equal(3, foo(1), "foo return check")
		s.Equal(4, foo(2), "foo return check")
		s.Equal(5, foo(1), "foo return check")
		s.Equal(6, foo(2), "foo return check")
		mock.Reset()

	})
}

// TestUnitSystemFuncApply 测试系统函数的mock
//  需要加上 -gcflags="-l"
func (s *MockerTestSuite) TestUnitSystemFuncApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		defer mock.Reset()

		mock.Func(rand.Int31).Return(int32(3))

		date, _ := time.Parse("2006-01-02 15:04:05", "2020-07-30 00:00:00")
		mock.Func(time.Now).Return(date)

		s.Equal(int32(3), rand.Int31(), "foo mock check")
		s.Equal(date, time.Now(), "foo mock check")

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
func (f *fake) Call2(i int) int {
	return i
}

//go:noinline
func (f *fake) call(i int) int {
	return i
}

// I 接口测试
type I interface {
	Call(int) int
	Call1(string) string
	call2(int32) int32
}

// I 接口实现1
type Impl1 struct {
}

func (i Impl1) Call(int) int {
	return 1
}

func (i Impl1) Call1(string) string {
	return "not ok"
}

func (i Impl1) call2(int32) int32 {
	return 1
}
