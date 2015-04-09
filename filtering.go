package main

import (
	"regexp"
	"strings"
)

type position struct {
	Start, End int
}

type match struct {
	str       string
	positions []position
}

func makeMatches(strs []string) []match {
	res := make([]match, len(strs))
	for i, s := range strs {
		res[i] = match{s, nil}
	}
	return res
}

func filtering(candidates []string, needle []rune) <-chan []match {
	result := make(chan []match)
	if len(needle) < 1 {
		go func() {
			res := makeMatches(candidates)
			result <- res
		}()
		return result
	}
	go func() {
		needles := strings.Split(strings.TrimRight(string(needle), " "), " ")
		res := []match{}
	candLoop:
		for _, cand := range candidates {
			m := match{cand, []position{}}
			restPos := 0
			for _, n := range needles {
				idx := strings.Index(cand[restPos:], n)
				if idx == -1 {
					continue candLoop
				}
				m.positions = append(m.positions, position{restPos + idx, restPos + idx + len(n)})
				restPos += idx + len(n)
			}
			res = append(res, m)
		}
		result <- res
	}()
	return result
}

func regexpFiltering(candidates []string, needle []rune) <-chan []match {
	result := make(chan []match)
	if len(needle) < 1 {
		go func() {
			res := makeMatches(candidates)
			result <- res
		}()
		return result
	}
	go func() {
		reg, err := regexp.Compile(string(needle))
		if err != nil {
			result <- []match{}
			return
		}
		res := make([]match, len(candidates))
		idx := 0
		for _, c := range candidates {
			pos := reg.FindStringIndex(c)
			if pos == nil {
				continue
			}
			res[idx] = match{str: c, positions: []position{position{pos[0], pos[1]}}}
			idx++
		}
		result <- res
	}()
	return result
}
