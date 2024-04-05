package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adamhu714/chirpy/internal/database"
)

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerGetChirps: Error connecting database")
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		log.Printf("handlerGetChirps: Error getting chirps")
		return
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("handlerGetChirps: Error marshalling json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
