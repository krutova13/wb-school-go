package main

import (
	"sort"
	"strings"
)

// FindAnagrams находит все множества анаграмм по заданному словарю
func FindAnagrams(words []string) map[string][]string {
	anagrams := make(map[string][]string)
	groups := make(map[string][]string)
	set := make(map[string]struct{})

	for _, word := range words {
		set[word] = struct{}{}
	}

	for word := range set {
		sorted := sortString(word)
		groups[sorted] = append(groups[sorted], word)
	}

	for _, group := range groups {
		if len(group) > 1 {
			anagrams[group[0]] = group
		}
	}

	return anagrams
}

func sortString(s string) string {
	chars := strings.Split(s, "")
	sort.Strings(chars)
	return strings.Join(chars, "")
}
