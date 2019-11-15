package bot

import (
	"strings"
)

func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " a an on the of to "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}
