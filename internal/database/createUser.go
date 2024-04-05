package database

func (db *DB) CreateUser(email string) error {
	newUser := User{
		Id:    -1,
		Email: email,
	}

	// load db into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	newUser.Id = len(dbStructure.Users) + 1

	// add chirp to memory object
	dbStructure.Users[newUser.Id] = newUser

	//write new db to file
	err = db.writeDB(dbStructure)

	return err
}
