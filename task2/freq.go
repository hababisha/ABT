package task2

import (
	"strings"
	"unicode"
)
func clean(s string) []string{
	cleanS := strings.ToLower(s)
	ss := ""
	for _, c := range cleanS{
		if unicode.IsLetter(c){
			ss += string(c)
		}
		if string(c) == " "{
			ss += " "
		}
	}

	
	return strings.Split(ss, " ")
}
func FreqCount(s string) map[string]int{	
	words := clean(s)
	m := make(map[string]int)

	for _, w := range words{
		m[w] += 1
	}
	return m
}