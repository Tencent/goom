// Package mocker_test 对 mocker 包的测试
// 当前文件实现了对 mocker.go 的单测
package mocker_test

import (
	"math/rand"
	"testing"
	"time"

	"git.code.oa.com/goom/mocker"
	"github.com/stretchr/testify/suite"
)

// TestUnitBuilderSysTestSuite 系统库 mock 测试入口
func TestUnitBuilderSysTestSuite(t *testing.T) {
	mocker.OpenDebug()
	suite.Run(t, new(mockerSysTestSuite))
}

type mockerSysTestSuite struct {
	suite.Suite
}

// TestUnitSystemFuncApply 测试系统函数的 mock
// 需要加上 -gcflags="-l"
// 在bazel 构建环境下, 因为系统库不支持开启 gcflags=-l ,所以暂不支持系统库中的短函数 mock
func (s *mockerSysTestSuite) TestUnitSystemFuncApply() {
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
