package ast_walker

import "fmt"

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

