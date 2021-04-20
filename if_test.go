// Package mocker_test 对mocker包的测试
// 当前文件实现了对if.go的单测
package mocker_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestUnitIfTestSuite(t *testing.T) {
	suite.Run(t, new(ifTestSuite))
}

type ifTestSuite struct {
	suite.Suite
}

func (s *ifTestSuite) TestIfAndReturn() {
	s.Run("success", func() {

	})
}
