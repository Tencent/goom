package mocker_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"git.code.oa.com/goom/mocker"
)

// TestUnitOriginTestSuite 测试入口
func TestUnitOriginTestSuite(t *testing.T) {
	suite.Run(t, new(OriginTestSuite))
}

// OriginTestSuite Builder测试套件
type OriginTestSuite struct {
	suite.Suite
}

// TestUnitFunc 测试调用原函数mock return
func (s *OriginTestSuite) TestCallOrigin() {
	s.Run("success", func() {
		mb := mocker.Create()

		// 定义原函数,用于占位,实际不会执行该函数体
		var origin = func(i int) int {
			fmt.Println("origin func placeholder")
			return 0 + i
		}

		mb.Func(fun1).Origin(&origin).Apply(func(i int) int {
			originResult := origin(i)
			return originResult + 100
		})

		s.Equal(101, fun1(1), "fun1 mock check")

		mb.Reset()

		s.Equal(1, fun1(1), "fun1 mock reset check")
	})
}
