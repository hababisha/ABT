package task2

import (
	"unicode"
)

func Palindrome(s string) bool{
	left, right := 0, len(s) -1


	
	for left < right{
		for !unicode.IsLetter(rune(s[left])){
			left += 1
		}
		for !unicode.IsLetter(rune(s[right])){
			right -= 1
		}

		if unicode.ToLower(rune(s[left])) != unicode.ToLower(rune(s[right])){
			return false
		}
		left += 1
		right -= 1
	}
	return true
}