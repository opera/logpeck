package logpeck

import (
	"testing"
)

func TestNewPeckFilter(*testing.T) {
	Keywords := ""
	filter := NewPeckFilter(Keywords)
	if len(filter.incl) != 0 || len(filter.excl) != 0 {
		panic(filter)
	}

	Keywords = "hello"
	filter = NewPeckFilter(Keywords)
	if len(filter.incl) != 1 || len(filter.excl) != 0 {
		panic(filter)
	}

	Keywords = "hello|^logpeck"
	filter = NewPeckFilter(Keywords)
	if len(filter.incl) != 1 || len(filter.excl) != 1 {
		panic(filter)
	}

	Keywords = "hello logpeck | golang|^ include "
	filter = NewPeckFilter(Keywords)
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
	Keywords := ""
	filter := NewPeckFilter(Keywords)

	if filter.Drop("hello") {
		panic(filter)
	}

	Keywords = "logpeck"
	filter = NewPeckFilter(Keywords)

	if !filter.Drop("hello") {
		panic(filter)
	}

	Keywords = "hello logpeck | golang|^ include "
	filter = NewPeckFilter(Keywords)

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
