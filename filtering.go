package main

import "strings"

type position struct {
	Start, End int
}

type match struct {
	str       string
	positions []position
}

func filtering(candidates []string, needle []rune) <-chan []match {
	result := make(chan []match)
	if len(needle) < 1 {
		go func() {
			res := make([]match, len(candidates))
			for i, s := range candidates {
				res[i] = match{s, nil}
			}
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
