package main

import (
	"slices"
	"strings"
)

func badWordsHandler(input string) string {
	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(input, " ")
	for i, word := range words {
		if slices.Contains(bad_words, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
