package mocker

import "git.code.oa.com/goom/mocker/internal/proxy"

// CachedMethodMocker 带缓存的方法Mocker
type CachedMethodMocker struct {
	*MethodMocker
	mCache  map[string]*MethodMocker
	umCache map[string]UnexportedMocker
}

func NewCachedMethodMocker(m *MethodMocker) *CachedMethodMocker {
	return &CachedMethodMocker{
		MethodMocker: m,
		mCache:       make(map[string]*MethodMocker, 16),
		umCache:      make(map[string]UnexportedMocker, 16),
	}
}

// CachedMethodMocker 设置结构体的方法名
func (m *CachedMethodMocker) Method(name string) ExportedMocker {
	if mocker, ok := m.mCache[name]; ok {
		return mocker
	}

	mocker := NewMethodMocker(m.pkgName, m.MethodMocker.structDef)
	mocker.Method(name)
	m.mCache[name] = mocker

	return mocker
}

// CachedMethodMocker 导出私有方法
func (m *CachedMethodMocker) ExportMethod(name string) UnexportedMocker {
	if mocker, ok := m.umCache[name]; ok {
		return mocker
	}

	mocker := NewMethodMocker(m.pkgName, m.MethodMocker.structDef)
	exportedMocker := mocker.ExportMethod(name)
	m.umCache[name] = exportedMocker

	return exportedMocker
}

// 清除mock
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

func NewCachedUnexportedMethodMocker(m *UnexportedMethodMocker) *CachedUnexportedMethodMocker {
	return &CachedUnexportedMethodMocker{
		UnexportedMethodMocker: m,
		mCache:                 make(map[string]*UnexportedMethodMocker, 16),
	}
}

// CachedMethodMocker 设置结构体的方法名
func (m *CachedUnexportedMethodMocker) Method(name string) UnexportedMocker {
	if mocker, ok := m.mCache[name]; ok {
		return mocker
	}

	mocker := NewUnexportedMethodMocker(m.pkgName, m.UnexportedMethodMocker.structName)
	mocker.Method(name)
	m.mCache[name] = mocker

	return mocker
}

// 清除mock
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

func NewCachedInterfaceMocker(interfaceMocker *DefaultInterfaceMocker) *CachedInterfaceMocker {
	return &CachedInterfaceMocker{
		DefaultInterfaceMocker: interfaceMocker,
		mCache:                 make(map[string]InterfaceMocker, 16),
		ctx:                    interfaceMocker.ctx,
	}
}

func (m *CachedInterfaceMocker) Method(name string) InterfaceMocker {
	if mocker, ok := m.mCache[name]; ok {
		return mocker
	}

	mocker := NewDefaultInterfaceMocker(m.pkgName, m.iface, m.ctx)
	mocker.Method(name)
	m.mCache[name] = mocker

	return mocker
}
