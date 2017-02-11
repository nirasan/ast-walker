package ast_walker

import (
	"fmt"
	"strings"
	"regexp"
)

type History struct {
	List []string
}

func NewHistory(l int) *History {
	return &History{ List: make([]string, l) }
}

func (h *History) Append(s string) *History {
	dst := NewHistory(len(h.List) + 1)
	copy(dst.List, h.List)
	dst.List[len(dst.List)-1] = s
	return dst
}

func (h *History) AppendN(s string, i int) *History {
	return h.Append(fmt.Sprintf("%s[%d]", s, i))
}

func (h *History) Path() string {
	if h == nil {
		return ""
	}
	return h.Join(".")
}

func (h *History) Join(sep string) string {
	if h == nil {
		return ""
	}
	return strings.Join(h.List, sep)
}

func (h *History) MatchString(pattern string) bool {
	if h == nil {
		return false
	}
	matched, err := regexp.MatchString(pattern, h.Path())
	if err != nil {
		return false
	}
	return matched
}

func (h *History) MatchRegex(re *regexp.Regexp) bool {
	if h == nil {
		return false
	}
	return re.MatchString(h.Path())
}