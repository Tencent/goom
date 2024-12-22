package mocker_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"git.woa.com/goom/mocker"
	"git.woa.com/goom/mocker/test"
	"github.com/stretchr/testify/suite"
)

// TestUnitUeVarTestSuite 测试入口
func TestUnitUeVarTestSuite(t *testing.T) {
	// 开启 debug
	// 1.可以查看 apply 和 reset 的状态日志
	// 2.查看 mock 调用日志
	mocker.OpenDebug()
	suite.Run(t, new(ueVarMockerTestSuite))
}

type ueVarMockerTestSuite struct {
	suite.Suite
	fakeErr error
}

func (s *ueVarMockerTestSuite) SetupTest() {
	s.fakeErr = errors.New("fake error")
}

func (s *ueVarMockerTestSuite) TestNewUeVarMock() {
	s.T().Log("args: ")
	for i := range os.Args {
		s.T().Log(os.Args[i], " ")
	}
	s.Run("success", func() {
		mocker := mocker.Create().UnExportedVar("git.woa.com/goom/mocker/test.unexportedGlobalIntVar", reflect.TypeOf(0))
		s.Equal(1, test.UnexportedGlobalIntVar(), "unexported global int var result check")
		mocker.Set(3)
		//fmt.Println(test.UnexportedGlobalIntVar())
		s.Equal(3, test.UnexportedGlobalIntVar(), "unexported global int var result check")
		mocker.Cancel()
		//fmt.Println(test.UnexportedGlobalIntVar())
		s.Equal(1, test.UnexportedGlobalIntVar(), "unexported global int var result check")
	})
}

func (s *ueVarMockerTestSuite) TestNewUeComplexVarMock() {
	testCases := []struct {
		path     string
		typ      reflect.Type
		initial  interface{}
		modified interface{}
		getter   func() interface{}
	}{
		{
			path:     "git.woa.com/goom/mocker/test.unexportedGlobalStrVar",
			typ:      reflect.TypeOf(""),
			initial:  "str",
			modified: "str1",
			getter:   func() interface{} { return test.UnexportedGlobalStrVar() },
		},
		{
			path:     "git.woa.com/goom/mocker/test.unexportedGlobalMapVar",
			typ:      reflect.TypeOf(map[string]int{}),
			initial:  map[string]int{"key": 1},
			modified: map[string]int{"key": 2},
			getter:   func() interface{} { return test.UnexportedGlobalMapVar() },
		},
		{
			path:     "git.woa.com/goom/mocker/test.unexportedGlobalArrVar",
			typ:      reflect.TypeOf([]int{}),
			initial:  []int{1, 2, 3},
			modified: []int{1, 2, 4},
			getter:   func() interface{} { return test.UnexportedGlobalArrVar() },
		},
	}

	for _, tc := range testCases {
		s.Run(tc.path, func() {
			m := mocker.Create().UnExportedVar(tc.path, tc.typ)
			s.Equal(tc.initial, tc.getter(), "unexported global var result check")
			m.Set(tc.modified)
			s.Equal(tc.modified, tc.getter(), "unexported global var result check")
			m.Cancel()
			s.Equal(tc.initial, tc.getter(), "unexported global var result check")
		})
	}
}

func (s *ueVarMockerTestSuite) TestNewUeConstMock() {
	s.T().Log("args: ")
	for i := range os.Args {
		s.T().Log(os.Args[i], " ")
	}
	s.Run("success", func() {
		mocker.Create().UnExportedVar("git.woa.com/goom/mocker/test.unexportedGlobalIntConst", reflect.TypeOf(1))
		fmt.Println("unexportedGlobalIntConst: ", test.UnexportedGlobalIntConst())
	})
}
