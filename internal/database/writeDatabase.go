package database

import (
	"encoding/json"
	"log"
	"os"
)

// db method creates new database file if needed
func (db *DB) endureDB() error {
	_, err := os.Stat(db.path)
	if err != nil {
		file, err := os.Create(db.path)
		if err != nil {
			file.Close()
			return err
		}
		file.Close()

		dbStructure := DBStructure{
			map[int]Chirp{},
			map[int]User{},
		}

		err = db.writeDB(dbStructure)
		if err != nil {
			log.Printf("Error in endureDB: %s", err.Error())
		}
	}
	return nil
}

// writes object in memory to database file
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	db.mux.RLock()
	err = os.WriteFile(db.path, data, 0744)
	db.mux.RUnlock()
	return err
}
