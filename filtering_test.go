package main

import (
	"reflect"
	"testing"
)

type result struct {
	C, F []string
	In   []rune
}

func TestFiltering(t *testing.T) {
	sets := []result{
		{
			[]string{"abca", "caba", "test"},
			[]string{"abca", "caba"},
			[]rune("ca"),
		},
		{
			[]string{"a.out", "b"},
			[]string{"a.out"},
			[]rune("a."),
		},
	}
	for _, r := range sets {
		res := filtering(r.C, r.In)
		if got := <-res; !reflect.DeepEqual(r.F, got) {
			t.Errorf("want %#v, but got %#v", r.F, got)
		}
	}
}
