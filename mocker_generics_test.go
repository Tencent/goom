//go:build go1.18
// +build go1.18

// Package mocker_test mock单元测试包
// 泛型相关的mock实现，特性支持参照来源于: https://taoshu.in/go/monkey/generic.html
package mocker_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tencent/goom"
)

type GT[T any] struct {
}

//go:noinline
func (gt *GT[T]) Hello() T {
	var t T
	fmt.Println("")
	return t
}

//go:noinline
func Hello[T int | string]() T {
	var t T
	fmt.Println("")
	return t
}

// TestUnitGenericsSuite 测试入口
func TestUnitGenericsSuite(t *testing.T) {
	// 开启 debug
	// 1.可以查看 apply 和 reset 的状态日志
	// 2.查看 mock 调用日志
	mocker.OpenDebug()
	suite.Run(t, new(mockerTestGenericsSuite))
}

type mockerTestGenericsSuite struct {
	suite.Suite
}

// TestGenericsMethod 测试泛型方法调用
func (s *mockerTestGenericsSuite) TestGenericsMethod() {
	s.Run("success", func() {
		myMocker := mocker.Create()
		defer myMocker.Reset()

		myMocker.Struct(&GT[string]{}).Method("Hello").
			Return("hello")
		myMocker.Struct(&GT[int]{}).Method("Hello").
			Return(1)

		var gt *GT[string]
		s.Equal("hello", gt.Hello(), "foo mock check")
		var gt1 *GT[int]
		s.Equal(1, gt1.Hello(), "foo mock check")
	})
}

// TestGenericsMethodFunc 测试泛型方法函数调用
func (s *mockerTestGenericsSuite) TestGenericsMethodFunc() {
	s.Run("success", func() {
		myMocker := mocker.Create()
		defer myMocker.Reset()

		hello1 := (&GT[string]{}).Hello
		myMocker.Func(hello1).Return("hello")
		hello2 := (&GT[int]{}).Hello
		myMocker.Func(hello2).Return(1)

		s.Equal("hello", hello1(), "foo mock check")
		s.Equal(1, hello2(), "foo mock check")
	})
}

// TestGenericsFunc 测试泛型函数调用
func (s *mockerTestGenericsSuite) TestGenericsFunc() {
	s.Run("success", func() {
		myMocker := mocker.Create()
		defer myMocker.Reset()

		myMocker.Func(Hello[string]).Return("hello")
		myMocker.Func(Hello[int]).Return(1)

		s.Equal("hello", Hello[string](), "foo mock check")
		s.Equal(1, Hello[int](), "foo mock check")
	})
}
