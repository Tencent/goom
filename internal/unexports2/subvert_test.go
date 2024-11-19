package unexports2

import (
	"fmt"
	"runtime"
	"testing"
)

type SubvertTester struct {
	A int
	a int
	int
}

func TestExposeFunction(t *testing.T) {
	if runtime.GOOS == "windows" {
		fmt.Printf("Skipping TestExposeFunction because it doesn't work in test binaries on this platform. Please run standalone_test.\n")
		return
	}

	assertDoesNotPanic(t, func() {
		exposed, err := ExposeFunction("reflect.valueMethodName", (func() string)(nil))
		if err != nil {
			t.Error(err)
			return
		}
		if exposed == nil {
			t.Errorf("exposed should not be nil")
			return
		}
		f := exposed.(func() string)
		expected := "git.woa.com/goom/mocker/internal/unexports2.getPanic"
		actual := f()
		if actual != expected {
			t.Errorf("Expected [%v] but got [%v]", expected, actual)
			return
		}
	})
}

func assertDoesNotPanic(t *testing.T, function func()) {
	if err := getPanic(function); err != nil {
		t.Errorf("Unexpected panic: %v", err)
	}
}

func getPanic(function func()) (result interface{}) {
	defer func() {
		if e := recover(); e != nil {
			result = e
		}
	}()

	function()
	return
}
