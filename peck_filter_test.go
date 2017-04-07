package logpeck

import (
	"testing"
)

func TestNewPeckFilter(*testing.T) {
	filterExpr := ""
	filter := NewPeckFilter(filterExpr)
	if len(filter.incl) != 0 || len(filter.excl) != 0 {
		panic(filter)
	}

	filterExpr = "hello"
	filter = NewPeckFilter(filterExpr)
	if len(filter.incl) != 1 || len(filter.excl) != 0 {
		panic(filter)
	}

	filterExpr = "hello|^logpeck"
	filter = NewPeckFilter(filterExpr)
	if len(filter.incl) != 1 || len(filter.excl) != 1 {
		panic(filter)
	}

	filterExpr = "hello logpeck | golang|^ include "
	filter = NewPeckFilter(filterExpr)
	if len(filter.incl) != 2 || len(filter.excl) != 1 {
		panic(filter)
	}
	if filter.incl[0] != "hello logpeck " ||
		filter.incl[1] != " golang" ||
		filter.excl[0] != " include " {
		panic(filter)
	}
}

func TestDrop(*testing.T) {
	filterExpr := "hello logpeck | golang|^ include "
	filter := NewPeckFilter(filterExpr)

	if !filter.Drop("hello") {
		panic(filter)
	}

	if !filter.Drop("hello logpeck include this") {
		panic(filter)
	}

	if filter.Drop("hello logpeck golang this") {
		panic(filter)
	}
}
