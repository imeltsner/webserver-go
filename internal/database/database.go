package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if dbData.Chirps == nil {
		dbData.Chirps = map[int]Chirp{}
	}

	id := len(dbData.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbData.Chirps[id] = chirp

	err = db.writeDB(dbData)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, len(dbData.Chirps))
	i := 0
	for _, v := range dbData.Chirps {
		chirps[i] = v
		i++
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if chirp, ok := dbData.Chirps[id]; ok {
		return chirp, nil
	}

	return Chirp{}, errors.New("chirp not found")
}

func (db *DB) CreateUser(email string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if len(dbData.Users) == 0 {
		dbData.Users = map[int]User{}
	}

	id := len(dbData.Users) + 1
	user := User{
		ID:    id,
		Email: email,
	}
	dbData.Users[id] = user

	err = db.writeDB(dbData)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(db.path)
		if err != nil {
			log.Printf("Error creating database file: %s", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbFile, err := os.ReadFile(db.path)
	if err != nil {
		log.Printf("Error reading DB file %s", err)
		return DBStructure{}, err
	}

	dbData := DBStructure{}
	if len(dbFile) != 0 {
		err = json.Unmarshal(dbFile, &dbData)
		if err != nil {
			log.Printf("Error unmarshaling DB file %s", err)
			return DBStructure{}, err
		}
	}

	return dbData, nil
}

func (db *DB) writeDB(dbData DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbData)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}
