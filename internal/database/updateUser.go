package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) UpdateUser(email string, password string, id int) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Printf("UpdateUser: Error during password hashing:%s", err.Error())
	}

	newUser := User{
		Id:       id,
		Email:    email,
		Password: hashed,
	}

	// load db into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	
	newUser.IsChirpyRed = dbStructure.Users[id].IsChirpyRed

	// add user to memory object
	dbStructure.Users[id] = newUser

	// write new db to file
	err = db.writeDB(dbStructure)

	return err
}
