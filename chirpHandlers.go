package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, 500, "error gettings chirps")
		return
	}

	respondWithJSON(w, 200, data)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	chirpID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Printf("Error parsing ID %s", err)
		respondWithError(w, 500, "error parsing ID")
		return
	}

	data, err := cfg.db.GetChirp(int(chirpID))
	if err != nil {
		log.Printf("Chirp not found %s", err)
		respondWithError(w, 404, "chirp not found")
		return
	}

	respondWithJSON(w, 200, data)
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
