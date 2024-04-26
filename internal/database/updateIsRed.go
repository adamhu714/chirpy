package database

func (db *DB) UpdateIsRed(id int, state bool) error {

	// load db into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	newUser := User{
		Id:          id,
		Email:       dbStructure.Users[id].Email,
		Password:    dbStructure.Users[id].Password,
		IsChirpyRed: state,
	}

	// add user to memory object
	dbStructure.Users[id] = newUser

	// write new db to file
	err = db.writeDB(dbStructure)

	return err
}
