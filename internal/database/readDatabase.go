package database

import (
	"encoding/json"
	"os"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// data structure to hold database in memory
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// db method loads database into memory data structure object
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	data, err := os.ReadFile(db.path)
	db.mux.RUnlock()
	if err != nil {
		return DBStructure{}, err
	}
	dbStructure := DBStructure{}
	err = json.Unmarshal(data, &dbStructure)
	return dbStructure, err
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	numOfChirp := len(dbStructure.Chirps)
	chirpSlice := make([]Chirp, numOfChirp)

	for i := 1; i < numOfChirp+1; i++ {
		chirpSlice[i-1] = dbStructure.Chirps[i]
	}

	return chirpSlice, nil
}
