package main

import (
	"regexp"
	"strings"
)

func filtering(candidates []string, needle []rune) <-chan []string {
	result := make(chan []string)
	if len(needle) < 1 {
		go func() {
			result <- candidates
		}()
		return result
	}
	go func() {
		needles := strings.Split(string(needle), " ")
		regs := []string{}
		for _, n := range needles {
			regs = append(regs, regexp.QuoteMeta(n))
		}
		matcher := regexp.MustCompile(strings.Join(regs, ".*"))
		res := []string{}
		for _, c := range candidates {
			if matcher.MatchString(c) {
				res = append(res, c)
			}
		}
		result <- res
	}()
	return result
}
