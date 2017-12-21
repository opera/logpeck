package logpeck

import (
	"strings"
)

type PeckFilter struct {
	incl      []string
	excl      []string
	have_incl bool
	have_excl bool
}

func NewPeckFilter(Keywords string) *PeckFilter {
	filter := &PeckFilter{have_incl: false, have_excl: false}
	substrs := strings.Split(Keywords, "|")
	for _, substr := range substrs {
		if substr == "" {
			continue
		}
		if substr[0] == '^' {
			filter.excl = append(filter.excl, substr[1:])
			filter.have_excl = true
		} else {
			filter.incl = append(filter.incl, substr)
			filter.have_incl = true
		}
	}
	return filter
}

func (p *PeckFilter) Drop(str string) bool {
	res := false
	for _, f := range p.incl {
		if strings.Contains(str, f) {
			res = false
			break
		} else {
			res = true
		}
	}
	if res {
		return true
	}
	for _, f := range p.excl {
		if strings.Contains(str, f) {
			return true
		}
	}
	return false
}
