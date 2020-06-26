package mocker_test

import (
	"git.code.oa.com/goom/mocker/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"

	"git.code.oa.com/goom/mocker"
)

// TestBuilderFunc 测试私有方法mock
func TestBuilderFunc(t *testing.T) {
	mb := mocker.Create("")
	mb.Func("fun1").Proxy(func(i int) int {
		return i * 3
	})

	mb.Func("fun2").Proxy(func(i int) int {
		return i * 3
	})

	assert.Equal(t, 3, fun1(1), "fun1 mock apply check")
	assert.Equal(t, 3, fun2(1), "fun2 mock apply check")

	mb.Reset()

	assert.Equal(t, 1, fun1(1), "fun1 mock cancel check")
	assert.Equal(t, 2, fun2(1), "fun1 mock cancel check")
}

// TestMockUnexportedMethod 测试结构体的私有方法mock
func TestBuilderStruct(t *testing.T) {
	mb := mocker.Create("")
	mb.Struct("fake").Method("call").Proxy(func(_ *fake, i int) int {
		return i * 2
	})

	f := &fake{}

	assert.Equal(t, 2, f.call(1), "call mock apply check")

	mb.Reset()

	assert.Equal(t, 1, f.call(1), "call mock cancel check")
}

// TestBuilderFuncDef 测试函数定义的mock
func TestBuilderFuncDef(t *testing.T) {
	logger.LogLevel = logger.DebugLevel
	logger.Log2Console(true)

	mb := mocker.Create("")
	mb.FuncDef((&fake{}).call).Proxy(func(_ *fake, i int) int {
		return i * 2
	})


	f := &fake{}

	assert.Equal(t, 2, f.call(1), "call mock apply check")

	mb.Reset()

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
