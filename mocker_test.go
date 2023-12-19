// Package mocker_test 对 mocker 包的测试
// 当前文件实现了对 mocker.go 的单测
package mocker_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	mocker "github.com/tencent/goom"
	"github.com/tencent/goom/test"
)

// TestUnitBuilderTestSuite 测试入口
func TestUnitBuilderTestSuite(t *testing.T) {
	// 开启 debug
	// 1.可以查看 apply 和 reset 的状态日志
	// 2.查看 mock 调用日志
	mocker.OpenDebug()
	suite.Run(t, new(mockerTestSuite))
}

type mockerTestSuite struct {
	suite.Suite
	fakeErr error
}

func (s *mockerTestSuite) SetupTest() {
	s.fakeErr = errors.New("fake error")
}

// TestUnitFuncApply 测试函数 mock apply
func (s *mockerTestSuite) TestUnitFuncApply() {
	s.T().Log("args: ")
	for i := range os.Args {
		s.T().Log(os.Args[i], " ")
	}
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Func(test.Foo).Apply(func(int) int {
			return 3
		})
		s.Equal(3, test.Foo(1), "test.Foo mock check")

		mock.Reset()
		s.Equal(1, test.Foo(1), "test.Foo mock reset check")
	})
}

// TestClosureFuncApply 测试闭包函数 mock apply
func (s *mockerTestSuite) TestClosureFuncApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		var r = 1
		mock.Func(test.Foo).Apply(func(int) int {
			return r
		})
		r = 3
		s.Equal(3, test.Foo(1), "test.Foo mock check")

		mock.Reset()
		s.Equal(1, test.Foo(1), "test.Foo mock reset check")
	})
}

// TestUnitFuncReturn 测试函数 mock return
func (s *mockerTestSuite) TestUnitFuncReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Func(test.Foo).When(1).Return(3)
		s.Equal(3, test.Foo(1), "test.Foo mock check")

		mock.Reset()
		s.Equal(1, test.Foo(1), "test.Foo mock reset check")
	})
}

// TestUnitUnexportedFuncApply 测试未导出函数 mock apply
func (s *mockerTestSuite) TestUnitUnexportedFuncApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Pkg("github.com/tencent/goom/test").ExportFunc("foo").Apply(func(i int) int {
			return i * 3
		})
		s.Equal(3, test.Invokefoo(1), "foo mock check")

		mock.Reset()
		s.Equal(1, test.Invokefoo(1), "foo mock reset check")
	})
}

// TestUnitUnexportedFuncReturn 测试未导出函数 mock return
func (s *mockerTestSuite) TestUnitUnexportedFuncReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Pkg("github.com/tencent/goom/test").ExportFunc("foo").As(func(i int) int {
			return i * 1
		}).Return(3)
		s.Equal(3, test.Invokefoo(1), "foo mock check")

		mock.Reset()
		s.Equal(1, test.Invokefoo(1), "foo mock reset check")
	})
}

// TestUnitMethodApply 测试结构体的方法 mock apply
func (s *mockerTestSuite) TestUnitMethodApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&test.Fake{}).Method("Call").Apply(func(*test.Fake, int) int {
			return 5
		})

		f := &test.Fake{}
		s.Equal(5, f.Call(1), "call mock check")

		mock.Reset()
		s.Equal(1, f.Call(1), "call mock reset check")
	})
}

// TestUnitMethodReturn 测试结构体的方法 mock return
func (s *mockerTestSuite) TestUnitMethodReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&test.Fake{}).Method("Call").Return(5).AndReturn(6)
		mock.Struct(&test.Fake{}).Method("Call2").Return(7).AndReturn(8)

		f := &test.Fake{}
		s.Equal(5, f.Call(1), "call mock check")
		s.Equal(6, f.Call(1), "call mock check")
		s.Equal(7, f.Call2(1), "call mock check")
		s.Equal(8, f.Call2(1), "call mock check")

		mock.Reset()
		s.Equal(1, f.Call(1), "call mock reset check")
	})
}

// TestUnitUnExportedMethodApply 测试结构体的未导出方法 mock apply
func (s *mockerTestSuite) TestUnitUnExportedMethodApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&test.Fake{}).ExportMethod("call").Apply(func(_ *test.Fake, i int) int {
			return i * 2
		})

		f := &test.Fake{}
		s.Equal(2, f.Invokecall(1), "call mock check")

		mock.Reset()
		s.Equal(1, f.Invokecall(1), "call mock reset check")
	})
}

// TestUnitUnexportedMethodReturn 测试结构体的未导出方法 mock return
func (s *mockerTestSuite) TestUnitUnexportedMethodReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Struct(&test.Fake{}).ExportMethod("call").As(func(_ *test.Fake, i int) int {
			return i * 2
		}).Return(6)

		f := &test.Fake{}
		s.Equal(6, f.Invokecall(1), "call mock check")

		mock.Reset()
		s.Equal(1, f.Invokecall(1), "call mock reset check")
	})
}

// TestUnitUnExportStruct 测试未导出结构体的方法 mock apply
func (s *mockerTestSuite) TestUnitUnExportStruct() {
	s.Run("success", func() {

		// _fake 从 test.fake 中拷贝过来
		type _fake struct {
			_ string // field1
			_ int    // field2
		}

		mock := mocker.Create()
		// 指定包名
		s.Equal("github.com/tencent/goom_test", mock.PkgName())

		mock.Pkg("github.com/tencent/goom/test").ExportStruct("*fake").
			Method("call").Apply(func(_ *_fake, i int) int {
			return i * 2
		})
		s.Equal("github.com/tencent/goom_test", mock.PkgName())

		f := test.NewUnexportedFake()
		s.Equal(2, f.Invokecall(1), "call mock check")

		mock.Reset()
		s.Equal(1, f.Invokecall(1), "call mock reset check")
	})
}

// TestMultiReturn 测试调用原函数多返回
func (s *mockerTestSuite) TestMultiReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Func(test.Foo).When(1).Return(3).AndReturn(2)
		s.Equal(3, test.Foo(1), "foo mock check")
		s.Equal(2, test.Foo(1), "foo mock check")

		mock.Reset()
		s.Equal(1, test.Foo(1), "foo mock reset check")
	})
}

// TestMultiReturns 测试调用原函数多返回
func (s *mockerTestSuite) TestMultiReturns() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Func(test.Foo).Returns(2, 3)
		s.Equal(2, test.Foo(1), "foo mock check")
		s.Equal(3, test.Foo(1), "foo mock check")

		mock.Func(test.Foo).Returns(4, 5)
		s.Equal(4, test.Foo(1), "foo mock check")
		s.Equal(5, test.Foo(1), "foo mock check")

		mock.Func(test.Foo).When(-1).Returns(6, 7)
		s.Equal(6, test.Foo(-1), "foo mock check")
		s.Equal(7, test.Foo(-1), "foo mock check")

		mock.Func(test.Foo).When(-2).Returns(8, 9)
		s.Equal(8, test.Foo(-2), "foo mock check")
		s.Equal(9, test.Foo(-2), "foo mock check")

		mock.Reset()
		s.Equal(1, test.Foo(1), "foo mock reset check")
	})
}

// TestUnitFuncTwiceApply 测试函数 mock apply 多次
func (s *mockerTestSuite) TestUnitFuncTwiceApply() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Func(test.Foo).When(1).Return(3)
		mock.Func(test.Foo).When(2).Return(6)
		s.Equal(3, test.Foo(1), "foo mock check")
		s.Equal(6, test.Foo(2), "foo mock check")

		mock.Reset()
		mock.Func(test.Foo).When(1).Return(2)
		s.Equal(2, test.Foo(1), "foo mock reset check")
		mock.Reset()
	})
}

// TestUnitDefaultReturn 测试函数 mock 返回默认值
func (s *mockerTestSuite) TestUnitDefaultReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		mock.Func(test.Foo).Return(3).AndReturn(4)
		mock.Func(test.Foo).Return(5).AndReturn(6)
		s.Equal(3, test.Foo(1), "foo return check")
		s.Equal(4, test.Foo(2), "foo return check")
		s.Equal(5, test.Foo(1), "foo return check")
		s.Equal(6, test.Foo(2), "foo return check")
		mock.Reset()
	})
}

// TestFakeReturn 测试返回 fake 值
func (s *mockerTestSuite) TestFakeReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		defer mock.Reset()

		mock.Func(test.Foo1).Return(&test.S1{
			Field1: "ok",
			Field2: 2,
		})
		s.Equal(&test.S{
			Field1: "ok",
			Field2: 2,
		}, test.Foo1(), "foo mock check")
	})
}

func (s *mockerTestSuite) TestUnitNilReturn() {
	s.Run("nil return", func() {
		mocker.Create().Func(test.GetS).Return(nil, s.fakeErr)
		res, err := test.GetS()
		s.Equal([]byte(nil), res)
		s.Equal(s.fakeErr, err)
	})
}

// TestVarMock 测试简单变量 mock
func (s *mockerTestSuite) TestVarMock() {
	s.Run("simple var mock", func() {
		mock := mocker.Create()
		mock.Var(&test.GlobalVar).Set(2)
		s.Equal(2, test.GlobalVar)
		mock.Reset()
		s.Equal(1, test.GlobalVar)
	})
}

// TestVarApply 测试变量应用 mock
func (s *mockerTestSuite) TestVarApply() {
	s.Run("var mock apply", func() {
		mock := mocker.Create()
		mock.Var(&test.GlobalVar).Apply(func() int {
			return 2
		})
		s.Equal(2, test.GlobalVar)
		mock.Reset()
		s.Equal(1, test.GlobalVar)
	})
}
