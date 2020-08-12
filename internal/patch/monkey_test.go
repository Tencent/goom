package patch_test

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"git.code.oa.com/goom/mocker/internal/logger"
	"git.code.oa.com/goom/mocker/internal/patch"

	"github.com/stretchr/testify/assert"
)

//go:noinline
func no() bool { return false }

//go:noinline
func yes() bool { return true }

func init() {
	logger.Log2Console(true)
}

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

func TestGC(t *testing.T) {
	value := true
	_, _ = patch.Patch(no, func() bool {
		return value
	})

	defer patch.UnpatchAll()
	runtime.GC()
	assert.True(t, no())
}

func TestSimple(t *testing.T) {
	assert.False(t, no())
	_, _ = patch.Patch(no, yes)
	assert.True(t, no())
	assert.True(t, patch.Unpatch(no))
	assert.False(t, no())
	assert.False(t, patch.Unpatch(no))
}

func TestGuard(t *testing.T) {
	var guard *patch.PatchGuard
	guard, _ = patch.Patch(no, func() bool {
		guard.Unpatch()
		defer guard.Restore()
		return !no()
	})

	for i := 0; i < 100; i++ {
		assert.True(t, no())
	}
	patch.Unpatch(no)
}

func TestUnpatchAll(t *testing.T) {
	assert.False(t, no())
	_, _ = patch.Patch(no, yes)
	assert.True(t, no())
	patch.UnpatchAll()
	assert.False(t, no())
}

type s struct{}

func (s *s) yes() bool { return true }

func TestWithInstanceMethod(t *testing.T) {
	i := &s{}

	assert.False(t, no())

	_, _ = patch.Patch(no, i.yes)

	assert.True(t, no())

	patch.Unpatch(no)
	assert.False(t, no())
}

type f struct{}

func (f *f) No() bool { return false }

func TestOnInstanceMethod(t *testing.T) {
	i := &f{}
	assert.False(t, i.No())
	_, _ = patch.PatchInstanceMethod(reflect.TypeOf(i), "No", func(_ *f) bool { return true })
	assert.True(t, i.No())
	assert.True(t, patch.UnpatchInstanceMethod(reflect.TypeOf(i), "No"))
	assert.False(t, i.No())
}

func TestNotFunction(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = patch.Patch(no, 1)
	})
	assert.Panics(t, func() {
		_, _ = patch.Patch(1, yes)
	})
}

func TestNotCompatible(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = patch.Patch(no, func() {})
	})
}
