package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hashicorp/go-set/v2"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	type response struct {
		Response string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 500, "error decoding JSON")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "chirp is too long")
		return
	}

	badWords := set.From[string]([]string{"kerfuffle", "sharbert", "fornax"})
	cleaned := removeBadWords(chirp.Body, *badWords)

	respondWithJSON(w, 200, response{
		Response: cleaned,
	})
}

func removeBadWords(words string, badWords set.Set[string]) string {
	wordArray := strings.Split(words, " ")
	for i, s := range wordArray {
		if badWords.Contains(strings.ToLower(s)) {
			wordArray[i] = "****"
		}
	}
	return strings.Join(wordArray, " ")
}
