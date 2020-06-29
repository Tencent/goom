# GOOM单测Mock框架
## 介绍
### 背景
1. 基于公司目前内部没有一款自己维护的适合公司业务迭代的mock框架，众多项目采用外界开源的gomonkey框架进行函数的mock，因其存在较一些的bug，不支持异包私有函数mock，同包私有方法mock等等问题, 加上团队目前实现一款改进版-无需unpath即可在mock过程中调用原函数的特性，可以支持到延迟模拟，参数更改，mock数据录制等功能，因此建立此项目
2. 目前有一半以上方案是基于gomock类似的实现方案, 此mock方案需要要求业务代码具备良好的接口设计，从而能顺利生成mock代码，而goom只需要指定函数名称或函数定义，就能支持到任意函数的mock，任意函数异常注入，延时模拟等扩展性功能

### 功能特性
1. 私有(未导出)函数(或方法)的mock, 普通函数的mock
2. mock过程中调用原函数(线程安全, 支持并发单测)
3. 异常注入，对函数调用支持异常注入，延迟模拟等稳定性测试
4. 所有操作都是并发安全的

### 将来
1. 支持数据驱动测试
2. 支持Mock锚点定义
3. 支持代码重构

## Install
```bash
go get git.code.oa.com/goom/mocker
```

## Example
```golang
// 在需要使用mock的测试文件import
import "git.code.oa.com/goom/mocker"
```
### 基本使用
#### 函数mock
```golang
// 函数定义如下
func foo(i int) int {
    return i
}

// mock示例
// 创建当前包的mocker
mock := mocker.Create()

// mock函数foo并设定返回值为1
mock.Func(foo).Return(1)

// mock函数foo并设置其代理函数
mock.Func(foo).Apply(func(int) int {
    return 1
})
```

#### 结构体方法mock
```golang
// 结构体定义如下
type fake struct{}

func (f *fake) Call(i int) int {
    return i
}

// 私有方法
func (f *fake) call(i int) int {
    return i
}

// mock示例
// 创建当前包的mocker
mock := mocker.Create()

// mock fake的方法Call并设置其代理函数
mock.Struct(&fake{}).Method("Call").Apply(func(_ *fake, i int) int {
    return i * 2
 })

// mock fake的方法Call并返回1
mock.Struct(&fake{}).Method("Call").Return(1)

// mock fake的私有方法call, mock前先调用ExportMethod将其导出，并设置其代理函数
mock.Struct(&fake{}).ExportMethod("call").Apply(func(_ *fake, i int) int {
    return i * 2
})

// mock fake的私有方法call, mock前先调用ExportMethod将其导出为函数类型，后续支持设置When, Return等
mock.Struct(&fake{}).ExportMethod("call").As(func(_ *fake, i int) int {
    return i * 2
}).Return(1)
```

### 高阶用法
#### 函数mock
```golang
// 针对其它包的mock示例
// 创建指定包的mocker，设置引用路径
mock := mocker.Package("git.code.oa.com/goom/mocker_test")

// mock函数foo1并设置其代理函数
mock.mb.ExportFunc("foo1").Apply(func(i int) int {
    return i * 3
})

// mock函数foo1并设置其返回值
mock.ExportFunc("foo1").As(func(i int) int {
    return 0
}).Return(1)
```

#### 结构体方法mock
```golang
// 针对其它包的mock示例
// 创建指定包的mocker，设置引用路径
mock := mocker.Package("git.code.oa.com/goom/mocker_test")

// mock其它包的私有结构体fake的私有方法call，并设置其代理函数
mock.ExportStruct("fake").Method("call").Apply(func(_ *fake, i int) int {
    return i * 2
})

// mock其它包的私有结构体fake的私有方法call，并设置其返回值
mock..ExportStruct("fake").Method("call").As(func(_ *fake, i int) int {
    return i * 2
}).Return(1)
```

#### 在代理函数中调用原函数
```golang
mb := mocker.Create()

// 定义原函数,用于占位,实际不会执行该函数体
var origin = func(i int) int {
    // 函数体长度必须大于一定值, 所以随意加一些代码进行填充
    fmt.Println("origin func placeholder")
    return 0 + i
}

mb.Func(foo1).Origin(&origin).Apply(func(i int) int {
    // 调用原函数
    originResult := origin(i)

    // 加入延时逻辑等
    time.Sleep(time.Seconds)

    return originResult + 100
})
```

## Contributor
@adrewchen、@miliao、@yongfuchen
