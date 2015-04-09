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
			[]string{"aa", "bb"},
			[]string{"aa", "bb"},
			[]rune(""),
		},
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
		{
			[]string{"README.md", "filtering.go", "filtering_test.go", "main.go"},
			[]string{"filtering.go", "filtering_test.go"},
			[]rune("fil go"),
		},
	}
	for _, r := range sets {
		res := filtering(r.C, r.In)
		if got := <-res; !reflect.DeepEqual(r.F, got) {
			t.Errorf("want %#v, but got %#v", r.F, got)
		}
	}
}
