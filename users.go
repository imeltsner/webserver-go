package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/imeltsner/webserver-go/internal/database"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error decoding user %s", err)
		respondWithError(w, 500, "error decoding user")
		return
	}

	user, err = cfg.db.CreateUser(user.Email)
	if err != nil {
		log.Printf("Error creating user %s", err)
		respondWithError(w, 500, "error creating user")
		return
	}

	respondWithJSON(w, 201, user)
}
