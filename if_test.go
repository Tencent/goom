package mocker_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestUnitIfTestSuite 测试入口
func TestUnitIfTestSuite(t *testing.T) {
	suite.Run(t, new(IfTestSuite))
}

// IfTestSuite If测试套件
type IfTestSuite struct {
	suite.Suite
}

// TestWhenAndReturn 多次返回不同的值
func (s *IfTestSuite) TestIfAndReturn() {
	s.Run("success", func() {

	})
}