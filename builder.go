package mocker

import (
	"fmt"
	"reflect"
	"strings"
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

// Struct 指定结构体名称
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (m *Builder) Struct(obj interface{}) *CachedMethodMocker {
	if mocker, ok := m.mCache[obj]; ok {
		return mocker.(*CachedMethodMocker)
	}

	mocker := &MethodMocker{
		baseMocker: newBaseMocker(m.pkgName),
		structDef:  obj,
	}

	cachedMocker := NewCachedMethodMocker(mocker)
	m.mockers = append(m.mockers, cachedMocker)
	m.mCache[obj] = cachedMocker

	return cachedMocker
}

// Func 指定函数定义
// funcdef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (m *Builder) Func(obj interface{}) *DefMocker {
	if mocker, ok := m.mCache[reflect.ValueOf(obj)]; ok {
		return mocker.(*DefMocker)
	}

	mocker := &DefMocker{
		baseMocker: newBaseMocker(m.pkgName),
		funcdef:    obj,
	}
	m.mockers = append(m.mockers, mocker)
	m.mCache[reflect.ValueOf(obj)] = mocker

	return mocker
}

// ExportStruct 导出私有结构体
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (m *Builder) ExportStruct(name string) *CachedUnexportedMethodMocker {
	if mocker, ok := m.mCache[m.pkgName+"_"+name]; ok {
		return mocker.(*CachedUnexportedMethodMocker)
	}

	structName := name

	if strings.Contains(name, "*") {
		structName = fmt.Sprintf("(%s)", name)
	}

	mocker := &UnexportedMethodMocker{
		baseMocker: newBaseMocker(m.pkgName),
		structName: structName,
	}

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
	if name == "" {
		panic("func name is empty")
	}

	if mocker, ok := m.mCache[m.pkgName+"_"+name]; ok {
		return mocker.(*UnexportedFuncMocker)
	}

	mocker := &UnexportedFuncMocker{
		baseMocker: newBaseMocker(m.pkgName),
		funcName:   name,
	}
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

// Create 创建Mock构建器
func Create() *Builder {
	return &Builder{pkgName: currentPackage(2), mCache: make(map[interface{}]interface{}, 30)}
}

// Package 创建Mock构建器
// Deprecated: 已支持在mock时设置pkg
func Package(_ string) *Builder {
	return &Builder{pkgName: currentPackage(2), mCache: make(map[interface{}]interface{}, 30)}
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

	mocker := m.MethodMocker.Method(name)
	m.mCache[name] = m.MethodMocker

	return mocker
}

// CachedMethodMocker 导出私有方法
func (m *CachedMethodMocker) ExportMethod(name string) UnexportedMocker {
	if mocker, ok := m.umCache[name]; ok {
		return mocker
	}

	mocker := m.MethodMocker.ExportMethod(name)
	m.umCache[name] = mocker

	return mocker
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

	mocker := m.UnexportedMethodMocker.Method(name)
	m.mCache[name] = m.UnexportedMethodMocker

	return mocker
}

// 清除mock
func (m *CachedUnexportedMethodMocker) Cancel() {
	for _, v := range m.mCache {
		v.Cancel()
	}
}
