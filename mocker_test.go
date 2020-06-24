package mocker_test

import (
	"git.code.oa.com/goom/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMocker(t *testing.T) {
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

//go:noinline
func fun1(i int) int {
	return i * 1
}

//go:noinline
func fun2(i int) int {
	return i * 2
}
