package mocker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.code.oa.com/goom/mocker"
)

// TestMockUnexportedFunc 测试私有方法mock
func TestMockUnexportedFunc(t *testing.T) {
	mb := mocker.Create("")
	mb.FuncName("fun1").Callback(func(i int) int {
		return i * 3
	}).Apply()

	mb.FuncName("fun2").Callback(func(i int) int {
		return i * 3
	}).Apply()

	assert.Equal(t, 3, fun1(1), "fun1 mock apply check")
	assert.Equal(t, 3, fun2(1), "fun2 mock apply check")

	mb.CancelAll()

	assert.Equal(t, 1, fun1(1), "fun1 mock cancel check")
	assert.Equal(t, 2, fun2(1), "fun1 mock cancel check")
}

// TestMockUnexportedMethod 测试类的私有方法mock
func TestMockUnexportedMethod(t *testing.T) {
	mb := mocker.Create("")
	mb.FuncName("(*fake).call").Callback(func(_ *fake, i int) int {
		return i * 2
	}).Apply()

	f := &fake{}

	assert.Equal(t, 2, f.call(1), "call mock apply check")

	mb.CancelAll()

	assert.Equal(t, 1, f.call(1), "call mock cancel check")
}

//go:noinline
func fun1(i int) int {
	return i * 1
}

//go:noinline
func fun2(i int) int {
	return i * 2
}

type fake struct{}

//go:noinline
func (f *fake) call(i int) int {
	return i
}
