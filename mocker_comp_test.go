// Package mocker_test 对 mocker 包的测试
// 当前文件实现了对 mocker_test.go,iface_test.go,debug_test.go 的单测对于不同 go 版本的兼容性测试
package mocker_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"git.woa.com/goom/mocker/test"
)

var versions = []string{
	"go1.16",
	"go1.17",
	"go1.18",
	"go1.19",
	"go1.20",
	"go1.21",
	"go1.22",
	"go1.23",
}

const testEnv = "MOCKER_COMPATIBILITY_TEST"

// TestCompatibility 测试针对不同 go 版本的兼容情况
func TestCompatibility(t *testing.T) {
	if os.Getenv(testEnv) == "true" {
		return
	}

	os.Setenv(testEnv, "true")
	for _, v := range versions {
		fmt.Printf("> [%s] start testing..\n", v)
		if err := test.Run(v, nil, "version"); err != nil {
			t.Errorf("[%s] env prepare fail: %v", v, err)
		}

		logHandler := func(log string) {
			if strings.Contains(log, "--- FAIL:") {
				t.Errorf("[%s] run fail: see details in the log above.", v)
			}
			fmt.Println(log)
		}
		vn := strings.Split(v, "go")[1]
		if err := test.Run(v, logHandler, "mod", "edit", "-go="+vn); err != nil {
			t.Errorf("[%s] run error: %v, see details in the log above.", v, err)
		}
		if err := test.Run(v, logHandler, "test", "-v", "-gcflags=\"all=-l\"", "-ldflags=\"-s=false\"", "."); err != nil {
			t.Errorf("[%s] run error: %v, see details in the log above.", v, err)
		}
		if t.Failed() {
			break
		}
		t.Logf("[%s] testing success.", v)
	}
}
