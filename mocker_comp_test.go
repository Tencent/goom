package mocker_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Jakegogo/goom_mocker/test"
)

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
		}
		if err := test.Run(v, logHandler, "test", "-v", "-gcflags=all=-l", "."); err != nil {
			t.Errorf("[%s] run error: %v, see details in the log above.", v, err)
		}
		if t.Failed() {
			break
		}
		t.Logf("[%s] testing success.", v)
	}
}
