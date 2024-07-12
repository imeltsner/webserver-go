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
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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
		dbData = DBStructure{
			Chirps: map[int]Chirp{},
		}
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

	// keys := make([]int, len(dbData.Chirps))
	// i := 0
	// for k := range dbData.Chirps {
	// 	keys[i] = k
	// 	i++
	// }
	// sort.Slice(keys, func(i, j int) bool {
	// 	return keys[i] < keys[j]
	// })

	// for _, k := range keys {
	// 	chirp, err := json.Marshal(dbData.Chirps[k])
	// 	if err != nil {
	// 		log.Printf("Error marshalling chirp %s", err)
	// 		return err
	// 	}
	// 	err = os.WriteFile(db.path, chirp, os.FileMode(0666))
	// 	if err != nil {
	// 		log.Printf("Error writing to DB file %s", err)
	// 		return err
	// 	}
	// }

	// return nil
}
