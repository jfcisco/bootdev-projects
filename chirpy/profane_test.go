package main

import (
	"fmt"
	"testing"
)

func TestRemoveProfaneWords(t *testing.T) {
	cases := []struct {
		sentence        string
		expectedCleaned string
	}{
		{
			sentence:        "This is a kerfuffle opinion I need to share with the world",
			expectedCleaned: "This is a **** opinion I need to share with the world",
		},
		{
			sentence:        "",
			expectedCleaned: "",
		},
		{
			sentence:        "Sharbert, SHARBERT, sharBerT",
			expectedCleaned: "****, ****, ****",
		},
		{
			sentence:        "Fornax should be censored, but f0rn4x should go through",
			expectedCleaned: "**** should be censored, but f0rn4x should go through",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			actual := removeProfaneWords(c.sentence)
			if actual != c.expectedCleaned {
				t.Errorf("RemoveProfaneWords() = \"%s\", want \"%s\"", actual, c.expectedCleaned)
			}
		})
	}
}
