package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/adamhu714/chirpy/internal/database"
)

func handlerGetChirpsId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("handlerGetChirpsId: Error while converting id")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerGetChirpsId: Error connecting database")
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		log.Printf("handlerGetChirpsId: Error getting chirps")
		return
	}

	if id < 1 || id > len(chirps) {
		log.Printf("handlerGetChirpsId: invalid chirp id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(chirps[id-1])
	if err != nil {
		log.Printf("handlerGetChirpsId: Error marshalling json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

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
