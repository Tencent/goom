//go:build !go1.18
// +build !go1.18

package patch

// IsGenericsFunc 是否为泛型函数
func IsGenericsFunc(_ string) bool {
	return false
}
