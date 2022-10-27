// Package patch_test patch 功能测试
package patch_test

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"git.woa.com/goom/mocker/internal/logger"
	"git.woa.com/goom/mocker/internal/patch"
	"git.woa.com/goom/mocker/internal/patch/test"

	"github.com/stretchr/testify/assert"
)

//init 初始化
func init() {
	logger.SetLog2Console(true)
}

// TestTimePatch timePatch 测试
func TestTimePatch(t *testing.T) {
	before := time.Now()

	guard, err := patch.Patch(time.Now, func() time.Time {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	})
	if err != nil {
		t.Error(err)
	}

	guard.Apply()

	during := time.Now()
	assert.True(t, patch.Unpatch(time.Now))
	after := time.Now()

	assert.Equal(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), during)
	assert.NotEqual(t, before, during)
	assert.NotEqual(t, during, after)
}

// TestGC GC 测试
func TestGC(t *testing.T) {
	value := true
	g, _ := patch.Patch(test.No, func() bool {
		return value
	})
	g.Apply()

	defer patch.UnpatchAll()
	runtime.GC()
	assert.True(t, test.No())
}

// TestSimple 测试降低函数
func TestSimple(t *testing.T) {
	assert.False(t, test.No())
	g, _ := patch.Patch(test.No, test.Yes)
	g.Apply()
	assert.True(t, test.No())
	assert.True(t, patch.Unpatch(test.No))
	assert.False(t, test.No())
	assert.False(t, patch.Unpatch(test.No))
}

// TestGuard 测试 guard.Apply()
func TestGuard(t *testing.T) {
	var guard *patch.Guard
	guard, _ = patch.Patch(test.No, func() bool {
		guard.Unpatch()
		defer guard.Restore()
		return !test.No()
	})
	guard.Apply()

	for i := 0; i < 100; i++ {
		assert.True(t, test.No())
	}
	patch.Unpatch(test.No)
}

//TestUnpatchAll 测试取消 patch
func TestUnpatchAll(t *testing.T) {
	assert.False(t, test.No())
	g, _ := patch.Patch(test.No, test.Yes)
	g.Apply()
	assert.True(t, test.No())
	patch.UnpatchAll()
	assert.False(t, test.No())
}

//TestWithInstanceMethod 测试实例方法
func TestWithInstanceMethod(t *testing.T) {
	i := &test.S{}

	assert.False(t, test.No())

	g, _ := patch.Patch(test.No, i.Yes)
	g.Apply()

	assert.True(t, test.No())

	patch.Unpatch(test.No)
	assert.False(t, test.No())
}

// TestOnInstanceMethod 测试实例方法
func TestOnInstanceMethod(t *testing.T) {
	i := &test.F{}
	assert.False(t, i.No())
	g, _ := patch.InstanceMethod(reflect.TypeOf(i), "No", func(_ *test.F) bool { return true })
	g.Apply()
	assert.True(t, i.No())
	assert.True(t, patch.UnpatchInstanceMethod(reflect.TypeOf(i), "No"))
	assert.False(t, i.No())
}

// TestNotFunction 测试 patch 到非函数类型
func TestNotFunction(t *testing.T) {
	assert.Panics(t, func() {
		g, _ := patch.Patch(test.No, 1)
		g.Apply()
	})
	assert.Panics(t, func() {
		g, _ := patch.Patch(1, test.Yes)
		g.Apply()
	})
}

// TestNotCompatible 测试 patch 到参数不一致的函数
func TestNotCompatible(t *testing.T) {
	assert.Panics(t, func() {
		g, _ := patch.Patch(test.No(), func() {})
		g.Apply()
	})
}
