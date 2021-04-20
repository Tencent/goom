// Package mocker_test 对mocker包的测试
// 当前文件实现了对iface.go的单测
package mocker_test

import (
	"testing"

	"git.code.oa.com/goom/mocker"
	"git.code.oa.com/goom/mocker/errobj"

	"github.com/stretchr/testify/suite"
)

// TestUnitIFaceTestSuite 接口Mock测试入口
func TestUnitIFaceTestSuite(t *testing.T) {
	suite.Run(t, new(ifaceMockerTestSuite))
}

type ifaceMockerTestSuite struct {
	suite.Suite
}

// TestUnitInterfaceApply 测试接口mock apply
func (s *ifaceMockerTestSuite) TestUnitInterfaceApply() {
	s.Run("success", func() {
		mock := mocker.Create()

		// 接口变量
		i := (I)(nil)

		// 将Mock应用到接口变量(仅对该变量有效)
		mock.Interface(&i).Method("Call").Apply(func(ctx *mocker.IContext, i int) int {
			return 3
		})
		mock.Interface(&i).Method("Call1").As(func(ctx *mocker.IContext, i string) string {
			return ""
		}).When("").Return("ok")

		s.Equal(3, i.Call(1), "interface mock check")
		s.Equal("ok", i.Call1(""), "interface mock check")

		s.NotNil(i, "interface var nil check")

		// Mock重置, 接口变量将恢复原来的值
		mock.Reset()

		s.Nil(i, "interface mock reset check")
	})
}

// TestUnitInterfaceReturn 测试接口mock return
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
		}).Return(int32(5))

		s.Equal(3, i.Call(1), "interface mock check")
		s.Equal("ok", i.Call1(""), "interface mock check")
		s.Equal(int32(5), i.call2(0), "interface mock check")

		s.NotNil(i, "interface var nil check")

		mock.Reset()

		s.Nil(i, "interface mock reset check")
	})
}

// TestUnitArgsNotMatch 测试接口mock参数不匹配情况
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

			// 将Mock应用到接口变量(仅对该变量有效)
			mock.Interface(&i).Method("Call").Apply(func(ctx *mocker.IContext) int {
				return 3
			})
		}()

		s.IsType(&errobj.IllegalParam{}, errobj.CauseOf(expectErr), "param check fail test")
	})
}

// I 接口测试
type I interface {
	Call(int) int
	Call1(string) string
	call2(int32) int32
}
