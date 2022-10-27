package patch

import (
	"testing"

	"git.woa.com/goom/mocker/internal/logger"
)

// nolint
//go:noinline
func Say() string {
	return "say"
}

// 测试地址修复功能
func Test_fixIns(t *testing.T) {
	logger.LogLevel = logger.InfoLevel

	originSay := func() string {
		println("just redundancy")
		return "anything"
	}
	guard, err := Trampoline(Say, func() string {
		return "mock " + originSay()
	}, &originSay)
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
