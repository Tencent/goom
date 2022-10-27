// Package mocker 定义了 mock 的外层用户使用 API 定义,
// 包括函数、方法、接口、未导出函数(或方法的)的 Mocker 的实现。
// 当前文件实现了当对同一方法或函数名进行重复构造时可以沿用缓存中已建好的 Mocker，
// 以防止在一个单测内重复构造 Mocker 时, 对上一个相同函数或方法的 Mocker 的内容规则造成覆盖。
package mocker

import (
	"strings"

	"git.woa.com/goom/mocker/internal/iface"
)

// CachedMethodMocker 带缓存的方法 Mocker,将同一个函数或方法的 Mocker 进行 cache
type CachedMethodMocker struct {
	*MethodMocker
	mCache  map[string]*MethodMocker
	umCache map[string]UnExportedMocker
}

// NewCachedMethodMocker 创建新的带缓存的方法 Mocker
func NewCachedMethodMocker(m *MethodMocker) *CachedMethodMocker {
	return &CachedMethodMocker{
		MethodMocker: m,
		mCache:       make(map[string]*MethodMocker, 16),
		umCache:      make(map[string]UnExportedMocker, 16),
	}
}

// String mock 的名称
func (m *CachedMethodMocker) String() string {
	s := make([]string, 0, len(m.mCache)+len(m.umCache))
	for _, v := range m.mCache {
		s = append(s, v.String())
	}
	for _, v := range m.umCache {
		s = append(s, v.String())
	}
	return strings.Join(s, ", ")
}

// Method 设置结构体的方法名
func (m *CachedMethodMocker) Method(name string) ExportedMocker {
	if mocker, ok := m.mCache[name]; ok && !mocker.Canceled() {
		return mocker
	}
	mocker := NewMethodMocker(m.pkgName, m.MethodMocker.structDef)
	mocker.Method(name)
	m.mCache[name] = mocker
	return mocker
}

// ExportMethod 导出私有方法
func (m *CachedMethodMocker) ExportMethod(name string) UnExportedMocker {
	if mocker, ok := m.umCache[name]; ok && !mocker.Canceled() {
		return mocker
	}
	mocker := NewMethodMocker(m.pkgName, m.MethodMocker.structDef)
	exportedMocker := mocker.ExportMethod(name)
	m.umCache[name] = exportedMocker
	return exportedMocker
}

// Cancel 取消 mock
func (m *CachedMethodMocker) Cancel() {
	for _, v := range m.mCache {
		v.Cancel()
	}
	for _, v := range m.umCache {
		v.Cancel()
	}
}

// CachedUnexportedMethodMocker 带缓存的未导出方法 Mocker
type CachedUnexportedMethodMocker struct {
	*UnexportedMethodMocker
	mockers map[string]*UnexportedMethodMocker
}

// NewCachedUnexportedMethodMocker 创建新的带缓存的未导出方法 Mocker
func NewCachedUnexportedMethodMocker(m *UnexportedMethodMocker) *CachedUnexportedMethodMocker {
	return &CachedUnexportedMethodMocker{
		UnexportedMethodMocker: m,
		mockers:                make(map[string]*UnexportedMethodMocker, 16),
	}
}

// String mock 的名称或描述
func (m *CachedUnexportedMethodMocker) String() string {
	s := make([]string, 0, len(m.mockers))
	for _, v := range m.mockers {
		s = append(s, v.String())
	}
	return strings.Join(s, ", ")
}

// Method 设置结构体的方法名
func (m *CachedUnexportedMethodMocker) Method(name string) UnExportedMocker {
	if mocker, ok := m.mockers[name]; ok && !mocker.Canceled() {
		return mocker
	}
	mocker := NewUnexportedMethodMocker(m.pkgName, m.UnexportedMethodMocker.structName)
	mocker.Method(name)
	m.mockers[name] = mocker
	return mocker
}

// Cancel 清除 mock
func (m *CachedUnexportedMethodMocker) Cancel() {
	for _, v := range m.mockers {
		v.Cancel()
	}
}

// CachedInterfaceMocker 带缓存的 Interface Mocker
type CachedInterfaceMocker struct {
	*DefaultInterfaceMocker
	mockers map[string]InterfaceMocker
	ctx     *iface.IContext
}

// NewCachedInterfaceMocker 创建新的带缓存的 Interface Mocker
func NewCachedInterfaceMocker(interfaceMocker *DefaultInterfaceMocker) *CachedInterfaceMocker {
	return &CachedInterfaceMocker{
		DefaultInterfaceMocker: interfaceMocker,
		mockers:                make(map[string]InterfaceMocker, 16),
		ctx:                    interfaceMocker.ctx,
	}
}

// String mock 的名称或描述
func (m *CachedInterfaceMocker) String() string {
	s := make([]string, 0, len(m.mockers))
	for _, v := range m.mockers {
		s = append(s, v.String())
	}
	return strings.Join(s, ", ")
}

// Method 指定方法名
func (m *CachedInterfaceMocker) Method(name string) InterfaceMocker {
	if mocker, ok := m.mockers[name]; ok && !mocker.Canceled() {
		return mocker
	}
	mocker := NewDefaultInterfaceMocker(m.pkgName, m.iFace, m.ctx)
	mocker.Method(name)
	m.mockers[name] = mocker
	return mocker
}

// Cancel 取消 mock
func (m *CachedInterfaceMocker) Cancel() {
	for _, v := range m.mockers {
		v.Cancel()
	}
}

// Canceled 是否取消了 mock
func (m *CachedInterfaceMocker) Canceled() bool {
	return m.ctx.Canceled()
}
