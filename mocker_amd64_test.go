// Package mocker_test 对 mocker 包的测试
// 当前文件实现了对 mocker.go 的单测
package mocker_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Jakegogo/goom_mocker/test"
	"github.com/stretchr/testify/suite"
)

// TestUnitAmd64TestSuite 测试入口
func TestUnitAmd64TestSuite(t *testing.T) {
	// 开启 debug
	// 1.可以查看 apply 和 reset 的状态日志
	// 2.查看 mock 调用日志
	mocker.OpenDebug()
	suite.Run(t, new(mockerTestAmd64Suite))
}

type mockerTestAmd64Suite struct {
	suite.Suite
	fakeErr error
}

func (s *mockerTestAmd64Suite) SetupTest() {
	s.fakeErr = errors.New("fake error")
}

// TestCallOrigin 测试调用原函数 mock return
func (s *mockerTestAmd64Suite) TestCallOrigin() {
	s.Run("success", func() {
		// 定义原函数,用于占位,实际不会执行该函数体
		var origin = func(i int) int {
			// 用于占位,实际不会执行该函数体, 但是必须编写
			fmt.Println("only for placeholder, will not call")
			// return 任意值
			return 0
		}

		mock := mocker.Create()
		mock.Func(test.Foo).Origin(&origin).Apply(func(i int) int {
			originResult := origin(i)
			fmt.Printf("arguments are %v\n", i)
			return originResult + 100
		})
		s.Equal(101, test.Foo(1), "foo mock check")

		mock.Reset()
		s.Equal(1, test.Foo(1), "foo mock reset check")
	})
}

// TestUnitSystemFuncApply 测试系统函数的 mock
// 需要加上 -gcflags="-l"
// 在bazel 构建环境下, 因为系统库不支持开启 gcflags=-l ,所以暂不支持系统库中的短函数 mock
func (s *mockerTestAmd64Suite) TestUnitSystemFuncApply() {
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

func (s *mockerTestAmd64Suite) TestUnitEmptyMatch() {
	s.Run("empty return", func() {
		mocker.Create().Func(time.Sleep).Return()
		time.Sleep(time.Second)
	})
}
