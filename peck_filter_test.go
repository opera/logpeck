package logpeck

import (
	"testing"
)

func TestNewPeckFilter(*testing.T) {
	keywords := ""
	filter := NewPeckFilter(keywords)
	if len(filter.incl) != 0 || len(filter.excl) != 0 {
		panic(filter)
	}

	keywords = "hello"
	filter = NewPeckFilter(keywords)
	if len(filter.incl) != 1 || len(filter.excl) != 0 {
		panic(filter)
	}

	keywords = "hello|^logpeck"
	filter = NewPeckFilter(keywords)
	if len(filter.incl) != 1 || len(filter.excl) != 1 {
		panic(filter)
	}

	keywords = "hello logpeck | golang|^ include "
	filter = NewPeckFilter(keywords)
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
	keywords := ""
	filter := NewPeckFilter(keywords)

	if filter.Drop("hello") {
		panic(filter)
	}

	keywords = "logpeck"
	filter = NewPeckFilter(keywords)

	if !filter.Drop("hello") {
		panic(filter)
	}

	keywords = "hello logpeck | golang|^ include "
	filter = NewPeckFilter(keywords)

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
