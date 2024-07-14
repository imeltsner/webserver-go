package database

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email, password string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if len(dbData.Users) == 0 {
		dbData.Users = map[int]DBUser{}
	}

	id := len(dbData.Users) + 1
	hashWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Printf("Error hashing password %s", err)
		return User{}, err
	}
	user := DBUser{
		ID:       id,
		Email:    email,
		Password: string(hashWord),
	}
	dbData.Users[id] = user

	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}

	res := User{
		ID:    user.ID,
		Email: user.Email,
	}

	return res, nil
}

func (db *DB) AuthenticateUser(user DBUser) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for k, v := range dbData.Users {
		if v.Email == user.Email {
			err = bcrypt.CompareHashAndPassword([]byte(v.Password), []byte(user.Password))
			if err != nil {
				return User{}, bcrypt.ErrMismatchedHashAndPassword
			} else {
				return User{ID: k, Email: v.Email}, nil
			}
		}
	}

	return User{}, errors.New("user not found")
}
