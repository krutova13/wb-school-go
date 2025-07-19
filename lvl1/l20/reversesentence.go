package l20

func ReverseSentence(s string) string {
	runes := []rune(s)

	reverse(runes, 0, len(runes)-1)

	start := 0
	for i := 0; i <= len(runes); i++ {
		if i == len(runes) || runes[i] == ' ' {
			reverse(runes, start, i-1)
			start = i + 1
		}
	}

	return string(runes)
}

func reverse(words []rune, start int, end int) {
	for start < end {
		words[start], words[end] = words[end], words[start]
		start++
		end--
	}
}
