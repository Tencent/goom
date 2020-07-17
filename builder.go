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
func (m *Builder) Struct(obj interface{}) *MethodMocker {
	if mocker, ok := m.mCache[obj]; ok {
		return mocker.(*MethodMocker)
	}

	mocker := &MethodMocker{
		baseMocker: newBaseMocker(m.pkgName),
		structDef:  obj,
	}
	m.mockers = append(m.mockers, mocker)
	m.mCache[obj] = mocker

	return mocker
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
func (m *Builder) ExportStruct(name string) *UnexportedMethodMocker {
	if mocker, ok := m.mCache[m.pkgName+"_"+name]; ok {
		return mocker.(*UnexportedMethodMocker)
	}

	structName := name

	if strings.Contains(name, "*") {
		structName = fmt.Sprintf("(%s)", name)
	}

	mocker := &UnexportedMethodMocker{
		baseMocker: newBaseMocker(m.pkgName),
		structName: structName,
	}
	m.mockers = append(m.mockers, mocker)
	m.mCache[m.pkgName+"_"+name] = mocker

	return mocker
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
