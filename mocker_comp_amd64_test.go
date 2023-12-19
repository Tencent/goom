// Package mocker_test 对 mocker 包的测试
// 当前文件实现了对 mocker_test.go,iface_test.go,debug_test.go 的单测对于不同 go 版本的兼容性测试
package mocker_test

var versions = []string{
	"go1.13",
	"go1.14",
	"go1.15",
	"go1.16",
	"go1.17",
	"go1.18",
	"go1.19",
	"go1.20",
}

const testEnv = "MOCKER_COMPATIBILITY_TEST"
