package mocker

import (
	"fmt"
	"runtime"
	"strings"
)

// MockerBuilder Mock构建器
type MockerBuilder struct {
	pkgname string
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
	return mocker
}

// FuncDef 指定函数定义
// funcdef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (m *MockerBuilder) FuncDef(funcdef interface{}) *Mocker {
	mocker := &Mocker{
		funcdef: funcdef,
	}
	return mocker
}

// Create 创建Mock构建起
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
	if i := strings.Index(callerName, ".(");i > -1 {
		return callerName[:i]
	}
	if i := strings.LastIndex(callerName, "."); i > -1 {
		return callerName[:i]
	}
	return callerName
}