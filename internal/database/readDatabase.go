package database

import (
	"encoding/json"
	"os"
	"time"
)

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

// data structure to hold database in memory
type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RevokedTokens map[string]time.Time `json:"revoked_tokens"`
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

func (db *DB) GetUsers() ([]User, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return []User{}, err
	}
	numOfUsers := len(dbStructure.Users)
	userSlice := make([]User, numOfUsers)

	for i := 1; i < numOfUsers+1; i++ {
		userSlice[i-1] = dbStructure.Users[i]
	}

	return userSlice, nil
}

func (db *DB) GetRevokedTokens() (map[string]time.Time, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return map[string]time.Time{}, err
	}

	return dbStructure.RevokedTokens, nil
}
