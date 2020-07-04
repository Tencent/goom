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

// TestWhenAndReturn 多次返回不同的值
func (s *IfTestSuite) TestIfAndReturn() {
	s.Run("success", func() {

	})
}
