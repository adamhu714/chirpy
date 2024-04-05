package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adamhu714/chirpy/internal/database"
)

func handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	email, err := validateUser(w, r)
	if err != nil {
		return
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerPostUsers - Error while connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.CreateUser(email)
	if err != nil {
		log.Printf("handlerPostUsers - Error while creating user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		log.Printf("Error adding user to database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respUser := users[len(users)-1]
	data, err := json.Marshal(respUser)
	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func validateUser(w http.ResponseWriter, r *http.Request) (string, error) {

	type requestParams struct {
		Email string `json:"email"`
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

	if len(requestBody.Email) == 0 {
		respBody := errorStruct{
			Error: "Email message not provided",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return "", errors.New("email message not provided")
	}

	return requestBody.Email, nil
}
