package logpeck

import (
	"strings"
)

// PeckFilter .
type PeckFilter struct {
	incl     []string
	excl     []string
	haveIncl bool
	haveExcl bool
}

// NewPeckFilter .
func NewPeckFilter(Keywords string) *PeckFilter {
	filter := &PeckFilter{haveIncl: false, haveExcl: false}
	substrs := strings.Split(Keywords, "|")
	for _, substr := range substrs {
		if substr == "" {
			continue
		}
		if substr[0] == '^' {
			filter.excl = append(filter.excl, substr[1:])
			filter.haveExcl = true
		} else {
			filter.incl = append(filter.incl, substr)
			filter.haveIncl = true
		}
	}
	return filter
}

// Drop return true if do not need processing
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
	SplitString("", "")
	return false
}
