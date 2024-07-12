package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hashicorp/go-set/v2"
	"github.com/imeltsner/webserver-go/internal/database"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := database.Chirp{}
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

	chirp, err = cfg.db.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, 500, "error creating chirp")
		return
	}

	respondWithJSON(w, 201, chirp)
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
