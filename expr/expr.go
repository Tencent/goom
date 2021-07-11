// Package expr 负责表达式执行, 执行结果用于When参数匹配(Matcher)
package expr

import (
	"reflect"
)

// Expr 表达式接口, 实现了equals、any、in、field(x)等表达式匹配
// 一般类型参数默认使用equals表达式, 实现了Expr接口的参数表达式将使用对于的表达式进行匹配参数
type Expr interface {
	// Eval 执行一个表达式,
	// 一般地, 如果执行结果是true, 则参数Match成功
	// input 表达式执行时的入参
	Eval(input []reflect.Value) (bool, error)
	// Resole 解析参数类型
	Resole(types []reflect.Type) error
}

// AnyExpr 和任意参数值比较
type AnyExpr struct {
	arg reflect.Value
}

func (a *AnyExpr) Resole(types []reflect.Type) error {
	return nil
}

// Eval 执行AnyExpr表达式
func (a *AnyExpr) Eval(_ []reflect.Value) (bool, error) {
	return true, nil
}

// EqualsExpr 表达式实现了两个参数是否相等的规则计算
type EqualsExpr struct {
	arg  interface{}
	argV reflect.Value
}

func (e *EqualsExpr) Resole(types []reflect.Type) error {
	e.argV = toValue(e.arg, types[0])
	return nil
}

// Eval 执行EqualsExpr表达式
func (e *EqualsExpr) Eval(input []reflect.Value) (bool, error) {
	if equal(e.argV, input[0]) {
		return true, nil
	}
	return false, nil
}

// InExpr 包含表达式执行
type InExpr struct {
	args  []interface{}
	exprs [][]Expr
}

func (i *InExpr) Resole(types []reflect.Type) error {
	exprs := make([][]Expr, 0)
	for _, v := range i.args {
		arg, ok := v.([]interface{})
		if !ok {
			arg = []interface{}{v}
		}
		// TODO results check
		expr := ToExpr(arg, types)
		exprs = append(exprs, expr)
	}
	i.exprs = exprs
	return nil
}

// InExpr InExpr 表达式执行
func (i *InExpr) Eval(input []reflect.Value) (bool, error) {
outer:
	for _, one := range i.exprs {
		if len(input) != len(one) {
			return false, nil
		}
		for i, arg := range one {
			v, err := arg.Eval([]reflect.Value{input[i]})
			if err != nil {
				return false, err
			}
			if !v {
				continue outer
			}
		}

		return true, nil
	}

	return false, nil
}
