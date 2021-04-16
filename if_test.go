// Package mocker_test 对mocker包的测试
// 当前文件实现了对if.go的单测
package mocker_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestUnitIfTestSuite 测试入口
func TestUnitIfTestSuite(t *testing.T) {
	suite.Run(t, new(IfTestSuite))
}

// IfTestSuite If测试套件
type IfTestSuite struct {
	suite.Suite
}

// TestIfAndReturn 多次返回不同的值
func (s *IfTestSuite) TestIfAndReturn() {
	s.Run("success", func() {

	})
}
