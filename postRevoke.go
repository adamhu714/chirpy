package main

import (
	"log"
	"net/http"

	"github.com/adamhu714/chirpy/internal/database"
)

func (cfg *apiConfig) handlerPostRevoke(w http.ResponseWriter, r *http.Request) {
	authHeaderContent := r.Header.Get("Authorization")
	if len(authHeaderContent) < 7 {
		log.Printf("handlerPostRevoke - bad jwt token provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerPostUsers - Error while connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.CreateRevokedToken(authHeaderContent[7:])
	if err != nil {
		log.Printf("Error retrieving users from database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
