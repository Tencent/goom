package mocker

import (
	"fmt"
	"reflect"
	"strings"

	"git.code.oa.com/goom/mocker/internal/proxy"
)

// Builder Mock构建器
type Builder struct {
	pkgName string
	mockers []Mocker

	mCache map[interface{}]interface{}
}

// Pkg 指定包名，当前包无需指定
func (m *Builder) Pkg(name string) *Builder {
	m.pkgName = name

	return m
}

// Create 创建Mock构建器
// 非线程安全的,不能在多协程中并发地mock或reset同一个函数
func Create() *Builder {
	return &Builder{pkgName: currentPackage(2), mCache: make(map[interface{}]interface{}, 30)}
}

// Interface 指定接口类型的变量定义
// iface 必须是指针类型, 比如 i为interface类型变量, iface传递&i
func (m *Builder) Interface(iface interface{}) *CachedInterfaceMocker {
	defer func() { m.pkgName = currentPackage(2) }()

	mKey := reflect.TypeOf(iface).String()
	if mocker, ok := m.mCache[mKey]; ok {
		return mocker.(*CachedInterfaceMocker)
	}

	// 创建InterfaceMocker
	// context和interface类型绑定
	mocker := NewDefaultInterfaceMocker(m.pkgName, iface, proxy.NewContext())

	cachedMocker := NewCachedInterfaceMocker(mocker)
	m.mockers = append(m.mockers, cachedMocker)
	m.mCache[mKey] = cachedMocker

	return cachedMocker
}

// Struct 指定结构体名称
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (m *Builder) Struct(obj interface{}) *CachedMethodMocker {
	defer func() { m.pkgName = currentPackage(2) }()

	mKey := reflect.ValueOf(obj).Type().String()
	if mocker, ok := m.mCache[mKey]; ok {
		return mocker.(*CachedMethodMocker)
	}

	mocker := NewMethodMocker(m.pkgName, obj)

	cachedMocker := NewCachedMethodMocker(mocker)
	m.mockers = append(m.mockers, cachedMocker)
	m.mCache[mKey] = cachedMocker

	return cachedMocker
}

// Func 指定函数定义
// funcdef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (m *Builder) Func(obj interface{}) *DefMocker {
	defer func() { m.pkgName = currentPackage(2) }()

	if mocker, ok := m.mCache[reflect.ValueOf(obj)]; ok {
		return mocker.(*DefMocker)
	}

	mocker := NewDefMocker(m.pkgName, obj)

	m.mockers = append(m.mockers, mocker)
	m.mCache[reflect.ValueOf(obj)] = mocker

	return mocker
}

// ExportStruct 导出私有结构体
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (m *Builder) ExportStruct(name string) *CachedUnexportedMethodMocker {
	defer func() { m.pkgName = currentPackage(2) }()

	if mocker, ok := m.mCache[m.pkgName+"_"+name]; ok {
		return mocker.(*CachedUnexportedMethodMocker)
	}

	structName := name

	if strings.Contains(name, "*") {
		structName = fmt.Sprintf("(%s)", name)
	}

	mocker := NewUnexportedMethodMocker(m.pkgName, structName)

	cachedMocker := NewCachedUnexportedMethodMocker(mocker)
	m.mockers = append(m.mockers, cachedMocker)
	m.mCache[m.pkgName+"_"+name] = cachedMocker

	return cachedMocker
}

// ExportFunc 导出私有函数
// 比如需要mock函数 foo()， 则name="pkgname.foo"
// 比如需要mock方法, pkgname.(*struct_name).method_name
// name string foo或者(*struct_name).method_name
func (m *Builder) ExportFunc(name string) *UnexportedFuncMocker {
	defer func() { m.pkgName = currentPackage(2) }()

	if name == "" {
		panic("func name is empty")
	}

	if mocker, ok := m.mCache[m.pkgName+"_"+name]; ok {
		return mocker.(*UnexportedFuncMocker)
	}

	mocker := NewUnexportedFuncMocker(m.pkgName, name)
	m.mockers = append(m.mockers, mocker)
	m.mCache[m.pkgName+"_"+name] = mocker

	return mocker
}

// Reset 取消当前builder的所有Mock
func (m *Builder) Reset() *Builder {
	for _, mocker := range m.mockers {
		mocker.Cancel()
	}

	return m
}

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
