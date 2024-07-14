package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/imeltsner/webserver-go/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := database.DBUser{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error decoding user %s", err)
		respondWithError(w, 500, "error decoding user")
		return
	}

	var newUser database.User
	newUser, err = cfg.db.CreateUser(user.Email, user.Password)
	if err != nil {
		log.Printf("Error creating user %s", err)
		respondWithError(w, 500, "error creating user")
		return
	}

	respondWithJSON(w, 201, newUser)
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	dbUser := database.DBUser{}
	err := decoder.Decode(&dbUser)
	if err != nil {
		log.Printf("Error decoding user %s", err)
		respondWithError(w, 500, "error decoding user")
		return
	}

	user, err := cfg.db.AuthenticateUser(dbUser)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		respondWithError(w, 401, "wrong password")
	} else if err != nil {
		respondWithError(w, 404, "user not found")
	} else {
		respondWithJSON(w, 200, user)
	}
}
