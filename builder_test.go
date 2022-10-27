package mocker

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(builderTestSuite))
}

type builderTestSuite struct {
	suite.Suite
}

func (s *builderTestSuite) TestBuilder_reset2CurPkg() {
	tests := []struct {
		name  string
		wants string
	}{{
		name:  "reset to current package",
		wants: "git.woa.com/goom/mocker",
	}}
	s.Run("success", func() {
		for _, tt := range tests {
			s.Run(tt.name, func() {
				b := Create()
				s.Equal(b.pkgName, tt.wants)
				b.Func(fake)
				s.Equal(b.pkgName, tt.wants)
			})
		}
	})
}

func (s *builderTestSuite) Test_currentPackage() {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "get current package name",
			want: "testing",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			// called by testing.tRunner()
			if got := currentPackage(); got != tt.want {
				if !s.Equal(got, tt.want) {
					debug.PrintStack()
				}
			}
		})
	}
}

func (s *builderTestSuite) Test_currentPkg() {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "skip 4",
			args: args{
				skip: 4,
			},
			want: "runtime",
		},
		{
			name: "skip 3",
			args: args{
				skip: 3,
			},
			want: "testing",
		},
		{
			name: "skip 2",
			args: args{
				skip: 2,
			},
			want: "github.com/stretchr/testify/suite",
		},
		{
			name: "skip 1(current package)",
			args: args{
				skip: 1,
			},
			want: "git.woa.com/goom/mocker",
		},
		{
			name: "skip 0(currentPkg define in package)",
			args: args{
				skip: 0,
			},
			want: "git.woa.com/goom/mocker",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			if got := currentPkg(tt.args.skip); got != tt.want {
				if !s.Equal(got, tt.want) {
					debug.PrintStack()
				}
			}
		})
	}
}

// fake 模拟函数
func fake() {}
