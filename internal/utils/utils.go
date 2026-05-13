package utils

import (
	"slices"
	"strings"
)

var profaneWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func RedactProfanity(text string) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
