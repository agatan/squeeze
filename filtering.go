package main

import (
	"regexp"
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
		reg := string(needle)
		matcher := regexp.MustCompile(regexp.QuoteMeta(reg))
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
