package main

import (
	"errors"
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

func matching(src string, needle []rune) (match, error) {
	switch currentMode {
	case normal:
		return matchingNormal(src, needle)
	case regex:
		return matchingRegex(src, needle)
	case fuzzy:
		return match{}, nil
	default:
		return match{}, nil
	}
}

func matchingNormal(src string, needle []rune) (match, error) {
	if len(needle) < 1 {
		return match{src, nil}, nil
	}
	needles := strings.Split(strings.TrimRight(string(needle), " "), " ")
	m := match{src, []position{}}
	restPos := 0
	for _, n := range needles {
		idx := strings.Index(src[restPos:], n)
		if idx == -1 {
			return m, errors.New("not match")
		}
		m.positions = append(m.positions, position{restPos + idx, restPos + idx + len(n)})
		restPos += idx + len(n)
	}
	return m, nil
}

func matchingRegex(src string, needle []rune) (match, error) {
	reg, err := regexp.Compile(string(needle))
	if err != nil {
		return match{}, errors.New("invalid regexp")
	}
	pos := reg.FindStringIndex(src)
	if pos == nil {
		return match{}, errors.New("not match")
	}
	return match{str: src, positions: []position{position{pos[0], pos[1]}}}, nil
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
		res := []match{}
		for _, cand := range candidates {
			m, err := matching(cand, needle)
			if err != nil {
				continue
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
