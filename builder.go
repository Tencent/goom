package mocker

import (
	"fmt"
	"runtime"
	"strings"
)

// Builder Mock构建器
type Builder struct {
	pkgname string
	mockers []Mocker
}

// Struct 指定结构体名称
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (m *Builder) Struct(name string) *MethodMocker {
	mocker := &MethodMocker{
		name: fmt.Sprintf("%s.%s", m.pkgname, name),
		namep:fmt.Sprintf("%s.(*%s)", m.pkgname, name)}
	m.mockers = append(m.mockers, mocker)

	return mocker
}

// Func 指定函数名称, 支持私有函数
// 比如需要mock函数 foo()， 则name="foo"
func (m *Builder) Func(name string) *FuncMocker {
	mocker := &FuncMocker{name: fmt.Sprintf("%s.%s", m.pkgname, name)}
	m.mockers = append(m.mockers, mocker)

	return mocker
}

// FuncDef 指定函数定义, 支持私有函数
// funcdef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (m *Builder) FuncDef(funcdef interface{}) *DefMocker {
	mocker := &DefMocker{funcdef: funcdef}
	m.mockers = append(m.mockers, mocker)

	return mocker
}

// Reset 取消package下的所有Mock, @see Builder.pkgname
func (m *Builder) Reset() *Builder {
	for _, mocker := range m.mockers {
		mocker.Cancel()
	}
	return m
}

// Create 创建Mock构建器
// pkgname string 包路径, 默认取当前包
func Create(pkgname string) *Builder {
	if pkgname == "" {
		pkgname = currentPackage(2)
	}
	return &Builder{
		pkgname: pkgname,
	}
}

// CurrentPackage 获取当前调用的包路径
func CurrentPackage() string {
	return currentPackage(2)
}

// currentPackage 获取调用者的包路径
func currentPackage(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	callerName := runtime.FuncForPC(pc).Name()
	if i := strings.Index(callerName, ".("); i > -1 {
		return callerName[:i]
	}
	if i := strings.LastIndex(callerName, "."); i > -1 {
		return callerName[:i]
	}
	return callerName
}
