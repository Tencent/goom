// Package mocker_test 对 mocker 包的测试
// 当前文件实现了对 iface.go 的单测
package mocker_test

import (
	"testing"

	"git.code.oa.com/goom/mocker"
	"git.code.oa.com/goom/mocker/erro"

	"github.com/stretchr/testify/suite"
)

// TestUnitIFaceTestSuite 接口 Mock 测试入口
func TestUnitIFaceTestSuite(t *testing.T) {
	mocker.OpenDebug()
	suite.Run(t, new(ifaceMockerTestSuite))
}

type ifaceMockerTestSuite struct {
	suite.Suite
}

// TestUnitInterfaceApply 测试接口 mock apply
func (s *ifaceMockerTestSuite) TestUnitInterfaceApply() {
	s.Run("success", func() {
		mock := mocker.Create()

		// 接口变量
		i := (I)(nil)

		// 将 Mock 应用到接口变量(仅对该变量有效)
		mock.Interface(&i).Method("Call").Apply(func(ctx *mocker.IContext, i int) int {
			return 3
		})
		mock.Interface(&i).Method("Call1").As(func(ctx *mocker.IContext, i string) string {
			return ""
		}).When("").Return("ok")

		s.Equal(3, i.Call(1), "interface mock check")
		s.Equal("ok", i.Call1(""), "interface mock check")

		s.NotNil(i, "interface var nil check")

		// Mock 重置, 接口变量将恢复原来的值
		mock.Reset()

		s.Nil(i, "interface mock reset check")
	})
}

// TestUnitInterfaceReturn 测试接口 mock return
func (s *ifaceMockerTestSuite) TestUnitInterfaceReturn() {
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
		}).Returns(int32(5), int32(6))

		s.Equal(3, i.Call(1), "interface mock check")
		s.Equal("ok", i.Call1(""), "interface mock check")
		s.Equal(int32(5), i.call2(0), "interface mock check")
		s.Equal(int32(6), i.call2(0), "interface mock check")

		s.NotNil(i, "interface var nil check")

		mock.Reset()

		s.Nil(i, "interface mock reset check")
	})
}

// TestUnitInterfaceTwice 测试多次接口 mock return
func (s *ifaceMockerTestSuite) TestUnitInterfaceAsTwice() {
	s.Run("success", func() {
		mock := mocker.Create()
		i := (I)(nil)

		mock.Interface(&i).Method("Call").As(func(ctx *mocker.IContext, i int) int {
			return 0
		}).When(1).Return(3)
		s.NotNil(i, "interface var nil check")
		s.Equal(3, i.Call(1), "interface mock check")
		mock.Reset()

		mock.Interface(&i).Method("Call").As(func(ctx *mocker.IContext, i int) int {
			return 0
		}).When(1).Return(4)
		s.NotNil(i, "interface var nil check")
		s.Equal(4, i.Call(1), "interface mock check")
		mock.Reset()

		s.Nil(i, "interface mock reset check")
	})
}

func (s *ifaceMockerTestSuite) TestUnitInterfaceApplyTwice() {
	s.Run("success", func() {
		mock := mocker.Create()
		i := (I)(nil)

		mock.Interface(&i).Method("Call").Apply(func(ctx *mocker.IContext, i int) int {
			return 1
		})
		mock.Interface(&i).Method("Call1").Apply(func(ctx *mocker.IContext, i string) string {
			return "1"
		})
		s.NotNil(i, "interface var nil check")
		s.Equal(1, i.Call(0), "interface mock check")
		s.Equal("1", i.Call1("0"), "interface mock check")
		mock.Reset()

		mock.Interface(&i).Method("Call").Apply(func(ctx *mocker.IContext, i int) int {
			return 2
		})
		mock.Interface(&i).Method("Call1").Apply(func(ctx *mocker.IContext, i string) string {
			return "2"
		})
		s.NotNil(i, "interface var nil check")
		s.Equal(2, i.Call(0), "interface mock check")
		s.Equal("2", i.Call1("0"), "interface mock check")
		mock.Reset()

		s.Nil(i, "interface mock reset check")
	})
}

// TestUnitArgsNotMatch 测试接口 mock 参数不匹配情况
func (s *ifaceMockerTestSuite) TestUnitArgsNotMatch() {
	s.Run("success", func() {

		var expectErr error
		func() {
			defer func() {
				if err := recover(); err != nil {
					expectErr, _ = err.(error)
				}
			}()
			mock := mocker.Create()

			// 接口变量
			i := (I)(nil)

			// 将 Mock 应用到接口变量(仅对该变量有效)
			mock.Interface(&i).Method("Call").Apply(func(ctx *mocker.IContext) int {
				return 3
			})
		}()

		s.IsType(&erro.IllegalParam{}, erro.CauseOf(expectErr), "param check fail test")
	})
}

// I 接口测试
type I interface {
	Call(int) int
	Call1(string) string
	call2(int32) int32
}
