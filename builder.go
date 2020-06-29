package mocker

import "fmt"

// Builder Mock构建器
type Builder struct {
	pkgname string
	mockers []Mocker
}

// Struct 指定结构体名称
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (m *Builder) Struct(obj interface{}) *MethodMocker {
	mocker := &MethodMocker{
		pkgname:    m.pkgname,
		baseMocker: newBackMocker(),
		structDef:  obj,
	}
	m.mockers = append(m.mockers, mocker)

	return mocker
}

// Func 指定函数定义
// funcdef 函数，比如 foo
// 方法的mock, 比如 &Struct{}.method
func (m *Builder) Func(obj interface{}) *DefMocker {
	mocker := &DefMocker{
		baseMocker: newBackMocker(),
		funcdef:    obj,
	}
	m.mockers = append(m.mockers, mocker)

	return mocker
}

// ExportStruct 导出私有结构体
// 比如需要mock结构体函数 (*conn).Write(b []byte)，则name="conn"
func (m *Builder) ExportStruct(name string) *UnexportedMethodMocker {
	mocker := &UnexportedMethodMocker{
		baseMocker: newBackMocker(),
		name:       fmt.Sprintf("%s.%s", m.pkgname, name),
		namep:      fmt.Sprintf("%s.(*%s)", m.pkgname, name),
	}
	m.mockers = append(m.mockers, mocker)

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

	mocker := &UnexportedFuncMocker{
		baseMocker: newBackMocker(),
		name:       fmt.Sprintf("%s.%s", m.pkgname, name)}
	m.mockers = append(m.mockers, mocker)

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
	pkgname := currentPackage(2)

	return &Builder{
		pkgname: pkgname,
	}
}

// Create 创建Mock构建器
// pkgname string 包路径,默认取当前包
func Package(pkgname string) *Builder {
	if pkgname == "" {
		pkgname = currentPackage(2)
	}

	return &Builder{
		pkgname: pkgname,
	}
}
