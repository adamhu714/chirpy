package database

func (db *DB) CreateChirp(body string) error {
	newChirp := Chirp{
		Id:   -1,
		Body: body,
	}

	// load db into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	newChirp.Id = len(dbStructure.Chirps) + 1

	// add chirp to memory object
	dbStructure.Chirps[newChirp.Id] = newChirp

	// write new db to file
	err = db.writeDB(dbStructure)

	return err
}
