package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email string, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Printf("CreateUser: Error during password hashing:%s", err.Error())
	}

	newUser := User{
		Id:          -1,
		Email:       email,
		Password:    hashed,
		IsChirpyRed: false,
	}

	// load db into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	newUser.Id = len(dbStructure.Users) + 1

	// add chirp to memory object
	dbStructure.Users[newUser.Id] = newUser

	// write new db to file
	err = db.writeDB(dbStructure)

	return err
}
