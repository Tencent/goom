package unexports

import (
	"strings"
)

// suggester 模糊匹配智能提示器
type suggester struct {
	key                                   string
	a, b                                  int
	suggestionA, suggestionB, suggestionC string
}

// newSuggester 创建提示器
func newSuggester(key string) *suggester {
	return &suggester{key: key}
}

// AddItem 添加匹配条目
func (s *suggester) AddItem(item string) {
	if fuzzyMatch(item, s.key, "/") {
		if s.b%3 == 0 {
			s.suggestionA = item
		} else if s.b%3 == 1 {
			s.suggestionB = item
		} else {
			s.suggestionC = item
		}
		s.b++
	} else if fuzzyMatch(item, s.key, ".") {
		if s.a%2 == 0 {
			s.suggestionB = item
		} else {
			s.suggestionC = item
		}
		s.a++
	}
}

// fuzzyMatch 模糊匹配,用于提供 suggestion
func fuzzyMatch(target, source, token string) bool {
	if len(target) == 0 || len(source) == 0 || len(token) == 0 {
		return false
	}
	keywords := strings.Split(source, token)
	keyword := keywords[len(keywords)-1]
	return strings.Contains(target, keyword)
}

// Suggestions 获取提示内容
func (s *suggester) Suggestions() []string {
	return []string{s.suggestionA, s.suggestionB, s.suggestionC}
}
