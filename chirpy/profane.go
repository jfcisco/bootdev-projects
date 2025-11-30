package main

import (
	"regexp"
	"strings"
)

var profaneWords []string = []string{"kerfuffle", "sharbert", "fornax"}

func removeProfaneWords(input string) string {
	// Split by word boundary (spaces, punctuation)
	words := regexp.MustCompile(`\b`).Split(input, -1)
	newWords := []string{}

	for _, word := range words {
		// Filter out profane words
		isProfane := false
		for _, profane := range profaneWords {
			// Use .EqualFold for case insensitive equals
			if strings.EqualFold(word, profane) {
				isProfane = true
				break
			}
		}

		if isProfane {
			newWords = append(newWords, "****")
		} else {
			newWords = append(newWords, word)
		}
	}
	return strings.Join(newWords, "")
}
