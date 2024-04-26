package database

func (db *DB) DeleteChirp(id int) error {

	// load db into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	newChirp := Chirp{
		Id:       dbStructure.Chirps[id].Id,
		AuthorId: dbStructure.Chirps[id].AuthorId,
		Body:     "Chirp has been deleted.",
	}

	// add chirp to memory object
	dbStructure.Chirps[newChirp.Id] = newChirp

	// write new db to file
	err = db.writeDB(dbStructure)

	return err
}
