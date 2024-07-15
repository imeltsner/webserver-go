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
	Chirps map[int]Chirp  `json:"chirps"`
	Users  map[int]DBUser `json:"users"`
}

type DBUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
