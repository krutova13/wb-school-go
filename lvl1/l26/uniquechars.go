package l26

import "strings"

func HasUniqueChars(s string) bool {
	s = strings.ToLower(s)
	chars := make(map[rune]struct{})
	for _, ch := range s {
		if _, exists := chars[ch]; exists {
			return false
		} else {
			chars[ch] = struct{}{}
		}
	}
	return true
}
