package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/adamhu714/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirpsId(w http.ResponseWriter, r *http.Request) {
	authorId, err := cfg.GetIDFromAccessToken(w, r)
	if err != nil {
		log.Printf("handlerDeleteChirpsId - Error getting id from token: %s", err.Error())
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("handlerDeleteChirpsId: Error while converting id")
	}

	// connect database
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerDeleteChirpsId: Error connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		log.Printf("handlerDeleteChirpsId: Error getting chirps")
		return
	}

	if id < 1 || id > len(chirps) {
		log.Printf("handlerDeleteChirpsId: invalid chirp id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if chirps[id-1].AuthorId != authorId {
		log.Printf("handlerDeleteChirpsId: wrong user attempting to delete chirp")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	db.DeleteChirp(id)

	w.WriteHeader(http.StatusOK)
}
