// Package mocker_test 对mocker包的测试
// 当前文件实现了对mocker.go的单测
package mocker_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"git.code.oa.com/goom/mocker"
	"git.code.oa.com/goom/mocker/testdata"
)

// TestUnitBuilderTestSuite 测试入口
func TestUnitBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(mockerTestSuite))
}

type mockerTestSuite struct {
	suite.Suite
	fakeErr error
}

func (s *mockerTestSuite) SetupTest() {
	s.fakeErr = errors.New("fake error")
}

// TestUnitFuncApply 测试函数mock apply
func (s *mockerTestSuite) TestUnitFuncApply() {
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
func (s *mockerTestSuite) TestUnitFuncReturn() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).When(1).Return(3)

		s.Equal(3, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestUnitUnexportedFuncApply 测试未导出函数mock apply
func (s *mockerTestSuite) TestUnitUnexportedFuncApply() {
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
func (s *mockerTestSuite) TestUnitUnexportedFuncReturn() {
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
func (s *mockerTestSuite) TestUnitMethodApply() {
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
func (s *mockerTestSuite) TestUnitMethodReturn() {
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

// TestUnitUnExportedMethodApply 测试结构体的未导出方法mock apply
func (s *mockerTestSuite) TestUnitUnExportedMethodApply() {
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
func (s *mockerTestSuite) TestUnitUnexportedMethodReturn() {
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

// TestUnitUnExportStruct 测试未导出结构体的方法mock apply
func (s *mockerTestSuite) TestUnitUnExportStruct() {
	s.Run("success", func() {
		// 指定包名
		mock := mocker.Create()
		s.Equal(mock.PkgName(), "git.code.oa.com/goom/mocker_test")

		mock.Pkg("git.code.oa.com/goom/mocker/testdata").ExportStruct("*Fake").
			Method("call").Apply(func(_ *fake, i int) int {
			return i * 2
		})
		s.Equal(mock.PkgName(), "git.code.oa.com/goom/mocker_test")

		f := &testdata.Fake{}

		s.Equal(2, f.Call(1), "call mock check")

		mock.Reset()

		s.Equal(1, f.Call(1), "call mock reset check")
	})
}

// TestCallOrigin 测试调用原函数mock return
func (s *mockerTestSuite) TestCallOrigin() {
	s.Run("success", func() {
		mock := mocker.Create()

		// 定义原函数,用于占位,实际不会执行该函数体
		var origin = func(i int) int {
			// 用于占位,实际不会执行该函数体, 但是必须编写
			fmt.Println("only for placeholder, will not call")
			// return 任意值
			return 0
		}

		mock.Func(foo).Origin(&origin).Apply(func(i int) int {
			originResult := origin(i)
			fmt.Printf("arguments are %v\n", i)
			return originResult + 100
		})

		s.Equal(101, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestMultiReturn 测试调用原函数多返回
func (s *mockerTestSuite) TestMultiReturn() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).When(1).Return(3).AndReturn(2)

		s.Equal(3, foo(1), "foo mock check")
		s.Equal(2, foo(1), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestMultiReturns 测试调用原函数多返回
func (s *mockerTestSuite) TestMultiReturns() {
	s.Run("success", func() {
		mock := mocker.Create()

		mock.Func(foo).Returns(2, 3)

		s.Equal(2, foo(1), "foo mock check")
		s.Equal(3, foo(1), "foo mock check")

		mock.Func(foo).Returns(4, 5)

		s.Equal(4, foo(1), "foo mock check")
		s.Equal(5, foo(1), "foo mock check")

		mock.Func(foo).When(-1).Returns(6, 7)

		s.Equal(6, foo(-1), "foo mock check")
		s.Equal(7, foo(-1), "foo mock check")

		mock.Func(foo).When(-2).Returns(8, 9)

		s.Equal(8, foo(-2), "foo mock check")
		s.Equal(9, foo(-2), "foo mock check")

		mock.Reset()

		s.Equal(1, foo(1), "foo mock reset check")
	})
}

// TestUnitFuncTwiceApply 测试函数mock apply多次
func (s *mockerTestSuite) TestUnitFuncTwiceApply() {
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
func (s *mockerTestSuite) TestUnitDefaultReturn() {
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
func (s *mockerTestSuite) TestUnitSystemFuncApply() {
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

// TestFakeReturn 测试返回fake值
func (s *mockerTestSuite) TestFakeReturn() {
	s.Run("success", func() {
		mock := mocker.Create()
		defer mock.Reset()

		mock.Func(foo1).Return(&S1{
			field1: "ok",
			field2: 2,
		})

		s.Equal(&S{
			field1: "ok",
			field2: 2,
		}, foo1(), "foo mock check")

	})
}

func (s *mockerTestSuite) TestUnitEmptyMatch() {
	s.Run("empty return", func() {
		mocker.Create().Func(time.Sleep).Return()
		time.Sleep(time.Second)
	})
}

func (s *mockerTestSuite) TestUnitNilReturn() {
	s.Run("nil return", func() {
		mocker.Create().Func(getS).Return(nil, s.fakeErr)

		res, err := getS()

		s.Equal([]byte(nil), res)
		s.Equal(s.fakeErr, err)
	})
}

// TestVarMock 测试简单变量mock
func (s *mockerTestSuite) TestVarMock() {
	s.Run("simple var mock", func() {
		mock := mocker.Create()
		mock.Var(&globalVar).Return(2)
		s.Equal(2, globalVar)
		mock.Reset()
		s.Equal(1, globalVar)
	})
}

// TestVarApply 测试变量应用mock
func (s *mockerTestSuite) TestVarApply() {
	s.Run("var mock apply", func() {
		mock := mocker.Create()
		mock.Var(&globalVar).Apply(func() int {
			return 2
		})
		s.Equal(2, globalVar)
		mock.Reset()
		s.Equal(1, globalVar)
	})
}

// globalVar 用于测试全局变量mock
var globalVar = 1

//go:noinline
func foo(i int) int {
	// check对defer的支持
	defer func() { fmt.Printf("defer\n") }()
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

type S struct {
	field1 string
	field2 int
}

type S1 struct {
	field1 string
	field2 int
}

//go:noinline
// foo1 foo1
func foo1() *S {
	return &S{
		field1: "ok",
		field2: 2,
	}
}

func getS() ([]byte, error) {
	return []byte("hello"), nil
}
