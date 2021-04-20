// Package mocker 定义了mock的外层用户使用API定义,
// 包括函数、方法、接口、未导出函数(或方法的)的Mocker的实现。
// 当前文件实现了Mocker接口各实现类的构造链创建，以便通过链式构造一个Mocker对象。
package mocker

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"git.code.oa.com/goom/mocker/internal/proxy"
)

// Builder Mock构建器, 负责创建一个链式构造器.
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

// PkgName 返回包名
func (b *Builder) PkgName() string {
	return b.pkgName
}

// Create 创建Mock构建器
// 非线程安全的,不能在多协程中并发地mock或reset同一个函数
func Create() *Builder {
	return &Builder{
		pkgName: currentPackage(),
		mCache:  make(map[interface{}]interface{}, 30),
	}
}

// Interface 指定接口类型的变量定义
// iFace 必须是指针类型, 比如 i为interface类型变量, iFace传递&i
func (b *Builder) Interface(iFace interface{}) *CachedInterfaceMocker {
	mKey := reflect.TypeOf(iFace).String()
	if mocker, ok := b.mCache[mKey]; ok {
		b.reset2CurPkg()

		return mocker.(*CachedInterfaceMocker)
	}

	// 创建InterfaceMocker
	// context和interface类型绑定
	mocker := NewDefaultInterfaceMocker(b.pkgName, iFace, proxy.NewContext())
	cachedMocker := NewCachedInterfaceMocker(mocker)

	b.cache(mKey, cachedMocker)
	b.reset2CurPkg()

	return cachedMocker
}

// cache 添加到缓存
func (b *Builder) cache(mKey interface{}, cachedMocker Mocker) {
	b.mockers = append(b.mockers, cachedMocker)
	b.mCache[mKey] = cachedMocker
}

// Struct 指定结构体名称
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (b *Builder) Struct(obj interface{}) *CachedMethodMocker {
	mKey := reflect.ValueOf(obj).Type().String()
	if mocker, ok := b.mCache[mKey]; ok {
		b.reset2CurPkg()
		return mocker.(*CachedMethodMocker)
	}

	mocker := NewMethodMocker(b.pkgName, obj)
	cachedMocker := NewCachedMethodMocker(mocker)

	b.cache(mKey, cachedMocker)
	b.reset2CurPkg()

	return cachedMocker
}

// Func 指定函数定义
// funcDef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (b *Builder) Func(obj interface{}) *DefMocker {
	if mocker, ok := b.mCache[reflect.ValueOf(obj)]; ok {
		b.reset2CurPkg()

		return mocker.(*DefMocker)
	}

	mocker := NewDefMocker(b.pkgName, obj)

	b.cache(reflect.ValueOf(obj), mocker)
	b.reset2CurPkg()

	return mocker
}

// ExportStruct 导出私有结构体
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (b *Builder) ExportStruct(name string) *CachedUnexportedMethodMocker {
	if mocker, ok := b.mCache[b.pkgName+"_"+name]; ok {
		b.reset2CurPkg()
		return mocker.(*CachedUnexportedMethodMocker)
	}

	structName := name
	if strings.Contains(name, "*") {
		structName = fmt.Sprintf("(%s)", name)
	}

	mocker := NewUnexportedMethodMocker(b.pkgName, structName)
	cachedMocker := NewCachedUnexportedMethodMocker(mocker)

	b.cache(b.pkgName+"_"+name, cachedMocker)
	b.reset2CurPkg()

	return cachedMocker
}

// ExportFunc 导出私有函数
// 比如需要mock函数 foo()， 则name="pkg_name.foo"
// 比如需要mock方法, pkg_name.(*struct_name).method_name
// name string foo或者(*struct_name).method_name
func (b *Builder) ExportFunc(name string) *UnexportedFuncMocker {
	if name == "" {
		panic("func name is empty")
	}

	if mocker, ok := b.mCache[b.pkgName+"_"+name]; ok {
		b.reset2CurPkg()
		return mocker.(*UnexportedFuncMocker)
	}

	mocker := NewUnexportedFuncMocker(b.pkgName, name)

	b.cache(b.pkgName+"_"+name, mocker)
	b.reset2CurPkg()

	return mocker
}

// Reset 取消当前builder的所有Mock
func (b *Builder) Reset() *Builder {
	for _, mocker := range b.mockers {
		mocker.Cancel()
	}

	return b
}

// reset2CurPkg 设置回当前的包
func (b *Builder) reset2CurPkg() {
	b.pkgName = currentPackage()
}

// currentPackageIndex 获取当前包的堆栈层次
const currentPackageIndex = 4

// currentPackage 获取当前调用的包路径
func currentPackage() string {
	return currentPkg(currentPackageIndex)
}

// currentPkg 获取调用者的包路径
func currentPkg(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	callerName := runtime.FuncForPC(pc).Name()

	if i := strings.Index(callerName, ".("); i > -1 {
		return callerName[:i]
	}

	if i := strings.LastIndex(callerName, "/"); i > -1 {
		realIndex := strings.Index(callerName[i:len(callerName)-1], ".")
		return callerName[:realIndex+i]
	}

	return callerName
}
