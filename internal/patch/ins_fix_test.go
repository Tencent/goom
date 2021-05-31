package patch

import (
	"git.code.oa.com/goom/mocker/internal/logger"
	"testing"
)

//go:noinline
func Say() string {
	return "say"
}

// 测试地址修复功能
func Test_fixIns(t *testing.T) {
	logger.LogLevel = logger.DebugLevel

	originfSay := func() string {
		println("just redundancy")
		return "anything"
	}
	guard, err := Trampoline(Say, func() string {
		return "mock " + originfSay()
	}, &originfSay)
	if err != nil {
		t.Errorf("patch error: %v", err)
	}
	guard.Apply()

	say := Say()
	if say != "mock say" {
		t.Fatalf("patch fatal, unexpected mock value returned: %v, expect: mock say", say)
	}

	guard.Unpatch()
	say = Say()
	if say != "say" {
		t.Errorf("unpatch fail: %v", say)
	}
}
