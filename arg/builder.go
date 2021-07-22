// Package arg 负责参数表达式构造和执行, 执行结果用于When参数匹配(Matcher)
package arg

// AnyValues 匹配任意参数值
var AnyValues = Any()

// Any 匹配任意参数值
func Any() *AnyExpr {
	return &AnyExpr{}
}

// Equals 创建参数比较表达式
func Equals(arg interface{}) *EqualsExpr {
	return &EqualsExpr{arg: arg}
}

// In 包含表达式的参数比较
func In(values ...interface{}) *InExpr {
	return &InExpr{
		args: values,
	}
}

// Field 属性值匹配表达式
func Field(name string) *Builder {
	return (&Builder{}).Field(name)
}

// Builder Expr表达式构建器, 根据规则构建Expr表达式子类对象
// TODO 实现表达式树
type Builder struct {
}

// Field 指定属性名称
func (b *Builder) Field(name string) *Builder {
	return b
}

// In 添加In字句
func (b *Builder) In(values ...interface{}) *Builder {
	return b
}
