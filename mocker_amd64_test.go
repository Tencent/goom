// Package mocker_test 对 mocker 包的测试
// 当前文件实现了对 mocker.go 的单测
package mocker_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	mocker "github.com/tencent/goom"
	"github.com/tencent/goom/test"
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

func (s *mockerTestAmd64Suite) TestUnitEmptyMatch() {
	s.Run("empty return", func() {
		mocker.Create().Func(time.Sleep).Return()
		time.Sleep(time.Second)
	})
}
