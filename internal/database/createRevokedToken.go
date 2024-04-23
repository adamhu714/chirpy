package database

import "time"

func (db *DB) CreateRevokedToken(token string) error {
	// load db into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// add chirp to memory object
	dbStructure.RevokedTokens[token] = time.Now().UTC()

	// write new db to file
	err = db.writeDB(dbStructure)

	return err
}
