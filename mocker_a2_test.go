// Package mocker_test 对mocker包的测试
// 当前文件实现了对mocker_test.go,iface_test.go,debug_test.go的单测对于不同go版本的兼容性测试
package mocker_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"git.code.oa.com/goom/mocker/internal/hack"
)

var versions = []string{
	"go1.11",
	"go1.12",
	"go1.13",
	"go1.14",
	"go1.15",
	"go1.16",
	"go1.17",
	"go1.18beta1",
}

const testEnv = "MOCKER_COMPATIBILITY_TEST"

// TestCompatibility 测试针对不同go版本的兼容情况
func TestCompatibility(t *testing.T) {
	if os.Getenv(testEnv) == "true" {
		return
	}
	os.Setenv(testEnv, "true")

	for _, v := range versions {
		fmt.Printf("> [%s] start testing..\n", v)
		if err := hack.Run(v, nil, "version"); err != nil {
			t.Fatalf("[%s] env prepare fail: %v", v, err)
		}

		fail := false
		logHandler := func(log string) {
			if strings.Contains(log, "--- FAIL:") {
				fail = true
			}
		}
		if err := hack.Run(v, logHandler, "test", "-v", "-gcflags=all=-l", "."); err != nil {
			t.Fatalf("[%s] run error: %v, see details in the log above.", v, err)
		}
		if fail {
			t.Fatalf("[%s] run fail: see details in the log above.", v)
		}
		t.Logf("[%s] testing success.", v)
	}
}
