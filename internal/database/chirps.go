package database

import "errors"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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

// GetAllChirps returns all chirps in the database
func (db *DB) GetAllChirps() ([]Chirp, error) {
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
