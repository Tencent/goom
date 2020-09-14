package mocker

import (
	"fmt"
	"reflect"
	"strings"

	"git.code.oa.com/goom/mocker/internal/proxy"
)

const currentPackageIndex = 2

// Builder Mock构建器
type Builder struct {
	pkgName string
	mockers []Mocker

	mCache map[interface{}]interface{}
}

// Pkg 指定包名，当前包无需指定
func (b *Builder) Pkg(name string) *Builder {
	b.pkgName = name

	return b
}

func (b *Builder) GetPkgName() string {
	return b.pkgName
}

// Create 创建Mock构建器
// 非线程安全的,不能在多协程中并发地mock或reset同一个函数
func Create() *Builder {
	return &Builder{pkgName: currentPackage(currentPackageIndex), mCache: make(map[interface{}]interface{}, 30)}
}

// Interface 指定接口类型的变量定义
// iface 必须是指针类型, 比如 i为interface类型变量, iface传递&i
func (b *Builder) Interface(iface interface{}) *CachedInterfaceMocker {
	mKey := reflect.TypeOf(iface).String()
	if mocker, ok := b.mCache[mKey]; ok {
		b.pkgName = currentPackage(currentPackageIndex)

		return mocker.(*CachedInterfaceMocker)
	}

	// 创建InterfaceMocker
	// context和interface类型绑定
	mocker := NewDefaultInterfaceMocker(b.pkgName, iface, proxy.NewContext())

	cachedMocker := NewCachedInterfaceMocker(mocker)
	b.mockers = append(b.mockers, cachedMocker)
	b.mCache[mKey] = cachedMocker

	b.pkgName = currentPackage(currentPackageIndex)

	return cachedMocker
}

// Struct 指定结构体名称
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (b *Builder) Struct(obj interface{}) *CachedMethodMocker {
	mKey := reflect.ValueOf(obj).Type().String()
	if mocker, ok := b.mCache[mKey]; ok {
		b.pkgName = currentPackage(currentPackageIndex)

		return mocker.(*CachedMethodMocker)
	}

	mocker := NewMethodMocker(b.pkgName, obj)

	cachedMocker := NewCachedMethodMocker(mocker)
	b.mockers = append(b.mockers, cachedMocker)
	b.mCache[mKey] = cachedMocker

	b.pkgName = currentPackage(currentPackageIndex)

	return cachedMocker
}

// Func 指定函数定义
// funcdef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (b *Builder) Func(obj interface{}) *DefMocker {
	if mocker, ok := b.mCache[reflect.ValueOf(obj)]; ok {
		b.pkgName = currentPackage(currentPackageIndex)

		return mocker.(*DefMocker)
	}

	mocker := NewDefMocker(b.pkgName, obj)

	b.mockers = append(b.mockers, mocker)
	b.mCache[reflect.ValueOf(obj)] = mocker

	b.pkgName = currentPackage(currentPackageIndex)

	return mocker
}

// ExportStruct 导出私有结构体
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (b *Builder) ExportStruct(name string) *CachedUnexportedMethodMocker {
	if mocker, ok := b.mCache[b.pkgName+"_"+name]; ok {
		b.pkgName = currentPackage(currentPackageIndex)

		return mocker.(*CachedUnexportedMethodMocker)
	}

	structName := name

	if strings.Contains(name, "*") {
		structName = fmt.Sprintf("(%s)", name)
	}

	mocker := NewUnexportedMethodMocker(b.pkgName, structName)

	cachedMocker := NewCachedUnexportedMethodMocker(mocker)
	b.mockers = append(b.mockers, cachedMocker)
	b.mCache[b.pkgName+"_"+name] = cachedMocker

	b.pkgName = currentPackage(currentPackageIndex)

	return cachedMocker
}

// ExportFunc 导出私有函数
// 比如需要mock函数 foo()， 则name="pkgname.foo"
// 比如需要mock方法, pkgname.(*struct_name).method_name
// name string foo或者(*struct_name).method_name
func (b *Builder) ExportFunc(name string) *UnexportedFuncMocker {
	if name == "" {
		panic("func name is empty")
	}

	if mocker, ok := b.mCache[b.pkgName+"_"+name]; ok {
		b.pkgName = currentPackage(currentPackageIndex)

		return mocker.(*UnexportedFuncMocker)
	}

	mocker := NewUnexportedFuncMocker(b.pkgName, name)
	b.mockers = append(b.mockers, mocker)
	b.mCache[b.pkgName+"_"+name] = mocker

	b.pkgName = currentPackage(currentPackageIndex)

	return mocker
}

// Reset 取消当前builder的所有Mock
func (b *Builder) Reset() *Builder {
	for _, mocker := range b.mockers {
		mocker.Cancel()
	}

	return b
}
