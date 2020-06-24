package mocker

import (
	"fmt"
	"runtime"
	"strings"
)

// MockerBuilder Mock构建器
type MockerBuilder struct {
	pkgname string
	mockers []*Mocker
}

// FuncName 指定函数名称, 支持私有函数
// 比如:
// 私有类型的方法: (*conn).Write
// 私有函数: package.funcname
func (m *MockerBuilder) FuncName(funcname string) *Mocker {
	fullFuncName := fmt.Sprintf("%s.%s", m.pkgname, funcname)
	mocker := &Mocker{
		funcname: fullFuncName,
	}
	m.mockers = append(m.mockers, mocker)
	return mocker
}

// FuncDef 指定函数定义
// funcdef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (m *MockerBuilder) FuncDef(funcdef interface{}) *Mocker {
	mocker := &Mocker{
		funcdef: funcdef,
	}
	m.mockers = append(m.mockers, mocker)
	return mocker
}

// ApplyAll 全部应用package下所有的Mock, @see MockerBuilder.pkgname
func (m *MockerBuilder) ApplyAll() *MockerBuilder {
	for _, mocker := range m.mockers {
		mocker.Apply()
	}
	return m
}

// CancelAll 取消package下的所有Mock, @see MockerBuilder.pkgname
func (m *MockerBuilder) CancelAll() *MockerBuilder {
	for _, mocker := range m.mockers {
		mocker.Cancel()
	}
	return m
}

// ReApplyAll 全部应用package下所有的Mock, @see MockerBuilder.pkgname
func (m *MockerBuilder) ReApplyAll() *MockerBuilder {
	for _, mocker := range m.mockers {
		mocker.ReApply()
	}
	return m
}

// Create 创建Mock构建器
// pkgname string 包路径, 默认取当前包
func Create(pkgname string) *MockerBuilder {
	if pkgname == "" {
		pkgname = currentPackage(2)
	}
	return &MockerBuilder{
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
