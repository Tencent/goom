//go:build go1.20
// +build go1.20

package patch

import (
	"strings"
)

// IsGenericsFunc 是否为泛型函数
func IsGenericsFunc(name string) bool {
	return strings.Contains(name, "[...]")
}
