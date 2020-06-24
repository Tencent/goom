package mocker

import (
	"git.code.oa.com/goom/mocker/internal/patch"
	"git.code.oa.com/goom/mocker/internal/proxy"
)

// Mocker 对函数或方法进行mock
// 能支持到私有函数、私有类型的方法的Mock
type Mocker struct {
	funcname string
	funcdef interface{}
	callback interface{}
	origin interface{}
	pkgname string

	guard *patch.PatchGuard
	err error
}

// Callback 指定mock执行的回调函数
// mock回调函数, 需要和mock模板函数的签名保持一致
// 方法的参数签名写法比如: func(s *Struct, arg1, arg2 type), 其中第一个参数必须是接收体类型
func (m *Mocker) Callback(callback interface{}) *Mocker {
	m.callback = callback
	return m
}

// Apply 应用Mock
func (m *Mocker) Apply() *Mocker {
	if m.funcdef != nil {
		m.guard, m.err = proxy.StaticProxyByFunc(m.funcdef, m.callback, m.origin)
		return m.apply()
	}
	if m.funcname != "" {
		m.guard, m.err = proxy.StaticProxyByName(m.funcname, m.callback, m.origin)
		return m.apply()
	}
	return m
}

func (m *Mocker) apply() *Mocker {
	if m.guard != nil {
		m.guard.Apply()
	}
	return m
}

// Cancel 取消Mock
func (m *Mocker) Cancel() *Mocker {
	if m.guard != nil {
		m.guard.UnpatchWithLock()
	}
	return m
}

// ReApply 重新应用Mock
func (m *Mocker) ReApply() *Mocker {
	if m.guard != nil {
		m.guard.Restore()
	}
	return m
}
