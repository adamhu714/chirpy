package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/adamhu714/chirpy/internal/database"
)

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
	authorId, err := cfg.GetIDFromAccessToken(w, r)
	if err != nil {
		log.Printf("handlerPostChirps - Error getting id from token: %s", err.Error())
		return
	}

	body, err := validateChirp(w, r)
	if err != nil {
		return
	}

	// connect database
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("Error connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.CreateChirp(body, authorId)
	if err != nil {
		log.Printf("Error adding chirp to database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		log.Printf("Error adding chirp to database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respChirp := chirps[len(chirps)-1]
	data, err := json.Marshal(respChirp)
	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func validateChirp(w http.ResponseWriter, r *http.Request) (string, error) {

	type requestParams struct {
		Body string `json:"body"`
	}
	type errorStruct struct {
		Error string `json:"error"`
	}

	var requestBody requestParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)

	if err != nil {
		log.Printf("Error while json decoding: %s", err.Error())
		respBody := errorStruct{
			Error: "Something went wrong",
		}
		respondWithJSON(w, http.StatusInternalServerError, respBody)
		return "", err
	}

	if len(requestBody.Body) == 0 {
		respBody := errorStruct{
			Error: "Chirp message not provided",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return "", errors.New("chirp message not provided")
	}

	if len(requestBody.Body) > 140 {
		respBody := errorStruct{
			Error: "Chirp is too long",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return "", errors.New("chirp is too long")
	}

	splitBody := strings.Split(requestBody.Body, " ")
	for i := 0; i < len(splitBody); i++ {
		if strings.ToLower(splitBody[i]) == "kerfuffle" {
			splitBody[i] = "****"
		}
		if strings.ToLower(splitBody[i]) == "sharbert" {
			splitBody[i] = "****"
		}
		if strings.ToLower(splitBody[i]) == "fornax" {
			splitBody[i] = "****"
		}
	}
	cleanedBody := strings.Join(splitBody, " ")

	return cleanedBody, nil
}

func respondWithJSON(w http.ResponseWriter, code int, respBody interface{}) {
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
