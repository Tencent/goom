// Package mocker定义了mock的外层用户使用API定义,
// 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件实现了当对同一方法或函数名进行重复构造时可以沿用缓存中已建好的Mocker，
// 以防止在一个单测内重复构造Mocker时, 对上一个相同函数或方法的Mocker的内容规则造成覆盖。
package mocker

import "git.code.oa.com/goom/mocker/internal/proxy"

// CachedMethodMocker 带缓存的方法Mocker,将同一个函数或方法的Mocker进行cache
type CachedMethodMocker struct {
	*MethodMocker
	mCache  map[string]*MethodMocker
	umCache map[string]UnExportedMocker
}

// NewCachedMethodMocker 创建新的带缓存的方法Mocker
func NewCachedMethodMocker(m *MethodMocker) *CachedMethodMocker {
	return &CachedMethodMocker{
		MethodMocker: m,
		mCache:       make(map[string]*MethodMocker, 16),
		umCache:      make(map[string]UnExportedMocker, 16),
	}
}

// Method 设置结构体的方法名
func (m *CachedMethodMocker) Method(name string) ExportedMocker {
	if mocker, ok := m.mCache[name]; ok {
		return mocker
	}

	mocker := NewMethodMocker(m.pkgName, m.MethodMocker.structDef)
	mocker.Method(name)
	m.mCache[name] = mocker

	return mocker
}

// ExportMethod 导出私有方法
func (m *CachedMethodMocker) ExportMethod(name string) UnExportedMocker {
	if mocker, ok := m.umCache[name]; ok {
		return mocker
	}

	mocker := NewMethodMocker(m.pkgName, m.MethodMocker.structDef)
	exportedMocker := mocker.ExportMethod(name)
	m.umCache[name] = exportedMocker

	return exportedMocker
}

// Cancel 取消mock
func (m *CachedMethodMocker) Cancel() {
	for _, v := range m.mCache {
		v.Cancel()
	}

	for _, v := range m.umCache {
		v.Cancel()
	}
}

// CachedUnexportedMethodMocker 带缓存的未导出方法Mocker
type CachedUnexportedMethodMocker struct {
	*UnexportedMethodMocker
	mCache map[string]*UnexportedMethodMocker
}

// NewCachedUnexportedMethodMocker 创建新的带缓存的未导出方法Mocker
func NewCachedUnexportedMethodMocker(m *UnexportedMethodMocker) *CachedUnexportedMethodMocker {
	return &CachedUnexportedMethodMocker{
		UnexportedMethodMocker: m,
		mCache:                 make(map[string]*UnexportedMethodMocker, 16),
	}
}

// Method 设置结构体的方法名
func (m *CachedUnexportedMethodMocker) Method(name string) UnExportedMocker {
	if mocker, ok := m.mCache[name]; ok {
		return mocker
	}

	mocker := NewUnexportedMethodMocker(m.pkgName, m.UnexportedMethodMocker.structName)
	mocker.Method(name)
	m.mCache[name] = mocker

	return mocker
}

// Cancel 清除mock
func (m *CachedUnexportedMethodMocker) Cancel() {
	for _, v := range m.mCache {
		v.Cancel()
	}
}

// CachedInterfaceMocker 带缓存的Interface Mocker
type CachedInterfaceMocker struct {
	*DefaultInterfaceMocker
	mCache map[string]InterfaceMocker
	ctx    *proxy.IContext
}

// NewCachedInterfaceMocker 创建新的带缓存的Interface Mocker
func NewCachedInterfaceMocker(interfaceMocker *DefaultInterfaceMocker) *CachedInterfaceMocker {
	return &CachedInterfaceMocker{
		DefaultInterfaceMocker: interfaceMocker,
		mCache:                 make(map[string]InterfaceMocker, 16),
		ctx:                    interfaceMocker.ctx,
	}
}

// Method 指定方法名
func (m *CachedInterfaceMocker) Method(name string) InterfaceMocker {
	if mocker, ok := m.mCache[name]; ok {
		return mocker
	}

	mocker := NewDefaultInterfaceMocker(m.pkgName, m.iFace, m.ctx)
	mocker.Method(name)
	m.mCache[name] = mocker

	return mocker
}

// Cancel 取消mock
func (m *CachedInterfaceMocker) Cancel() {
	for _, v := range m.mCache {
		v.Cancel()
	}
}
