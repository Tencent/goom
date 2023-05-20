// Package test 被测对象都放在这个包
package test

import "fmt"

var toggle = false

// No 返回 false 的函数
//
//go:noinline
func No() bool {
	if toggle {
		fmt.Println("false")
	}
	return false
}

// Yes 返回 true 的函数
//
//go:noinline
func Yes() bool { return true }

// S 结构体
type S struct{}

// Yes 返回 true 的方法
func (s *S) Yes() bool { return true }

// F 结构体
type F struct{}

// No 返回 false 的方法
func (f *F) No() bool { return false }
