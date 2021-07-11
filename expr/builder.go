package expr

// ExprBuilder Expr表达式构建器, 根据规则构建Expr表达式子类对象
type ExprBuilder struct {
}

// Equals 创建参数比较表达式
func Equals(arg interface{}) *EqualsExpr {
	return &EqualsExpr{arg: arg}
}

// Any 和参数任意值比较
func Any() *AnyExpr {
	return &AnyExpr{}
}

// In 包含表达式的参数比较
func In(values []interface{}) *InExpr {
	return &InExpr{
		args: values,
	}
}
