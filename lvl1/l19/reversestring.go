package l19

func ReverseString(s string) string {
	runes := []rune(s)
	l := len(runes)
	reversed := make([]rune, l)
	for i, r := range runes {
		reversed[l-i-1] = r
	}
	return string(reversed)
}
