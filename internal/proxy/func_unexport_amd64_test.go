// Package proxy_test 对 proxy 包的测试
package proxy_test

import (
	"fmt"
	"testing"

	"git.woa.com/goom/mocker/internal/proxy"
)

// TestPrintMock 测试 mock fmt.Print
// 在bazel 构建环境下, 因为系统库不支持开启 gcflags=-l ,所以暂不支持系统库中的短函数 mock
func TestPrintMock(t *testing.T) {
	var trampoline = func(a ...interface{}) (n int, err error) {
		return 0, nil
	}

	// 静态代理函数
	patch, err := proxy.FuncName("fmt.Print", func(a ...interface{}) (n int, err error) {
		// 调用原来的函数
		return fmt.Println("called fmt.Print, args:", a)
	}, &trampoline)
	if err != nil {
		t.Error("mock print err:", err)
	}

	fmt.Println("ok", "1")
	patch.Unpatch()
	fmt.Println("unpatched")
	fmt.Println("ok", "2")
}
