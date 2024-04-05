package database

import (
	"sync"
)

// object to read and write from database file
type DB struct {
	path string
	mux  *sync.RWMutex
}

// new database connection struct
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.endureDB()
	return db, err
}
