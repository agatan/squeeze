package main

import (
	"reflect"
	"testing"
)

type result struct {
	Candidate []string
	Out       []match
	In        []rune
}

func makeTestMatch(str string, pos ...position) match {
	return match{str, pos}
}

var testCandidates = []string{
	"README.md",
	"filtering.go",
	"filtering_test.go",
	"main.go",
	"screen.go",
}

func TestRegexpFiltering(t *testing.T) {
	tests := []result{
		{
			testCandidates,
			[]match{
				makeTestMatch("README.md", position{0, len("README.md")}),
				makeTestMatch("filtering.go", position{0, len("filtering.go")}),
				makeTestMatch("filtering_test.go", position{0, len("filtering_test.go")}),
				makeTestMatch("main.go", position{0, len("main.go")}),
				makeTestMatch("screen.go", position{0, len("screen.go")}),
			},
			[]rune(".*"),
		},
	}
	for _, r := range tests {
		res := regexpFiltering(r.Candidate, r.In)
		if got := <-res; !reflect.DeepEqual(r.Out, got) {
			t.Errorf("regexpFiltering(%#v, `%s`) = %#v, want %#v", r.Candidate, string(r.In), got, r.Out)
		}
	}
}

func TestFiltering(t *testing.T) {
	sets := []result{
		{
			[]string{"aa", "bb"},
			[]match{match{"aa", nil}, match{"bb", nil}},
			[]rune(""),
		},
		{
			[]string{"abca", "caba", "test"},
			[]match{match{"abca", []position{position{2, 4}}}, match{"caba", []position{position{0, 2}}}},
			[]rune("ca"),
		},
		{
			[]string{"a.out", "b"},
			[]match{match{"a.out", []position{position{0, 2}}}},
			[]rune("a."),
		},
		{
			[]string{"README.md", "filtering.go", "filtering_test.go", "main.go"},
			[]match{match{"filtering.go", []position{position{0, 3}, position{10, 12}}}, match{"filtering_test.go", []position{position{0, 3}, position{15, 17}}}},
			[]rune("fil go"),
		},
		{
			[]string{"README.md", "filtering.go", "filtering_test.go"},
			[]match{makeTestMatch("filtering.go", position{0, 3}), makeTestMatch("filtering_test.go", position{0, 3})},
			[]rune("fil "),
		},
		{
			[]string{"README.md", "filtering.go", "filtering_test.go", "main.go"},
			[]match{},
			[]rune("go fil"),
		},
		{
			[]string{"README.md", "filtering.go", "filtering_test.go", "main.go", "go.main"},
			[]match{makeTestMatch("main.go", position{0, 4}, position{5, 7})},
			[]rune("main go"),
		},
	}
	for _, r := range sets {
		res := filtering(r.Candidate, r.In)
		if got := <-res; !reflect.DeepEqual(r.Out, got) {
			t.Errorf("filtering(%#v, `%s`) = %#v, want %#v", r.Candidate, string(r.In), got, r.Out)
		}
	}
}
