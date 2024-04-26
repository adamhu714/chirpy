package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adamhu714/chirpy/internal/database"
)

func handlerPostPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	userId, err := validatePostPolkaWebhook(w, r)
	if err != nil {
		return
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerPostUsers - Error while connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		log.Printf("handlerPostPolkaWebhooks - error getting users: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userId < 1 || userId > len(users) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = db.UpdateIsRed(userId, true)
	if err != nil {
		log.Printf("handlerPostUsers - Error while creating user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validatePostPolkaWebhook(w http.ResponseWriter, r *http.Request) (int, error) {
	type requestParams struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
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
		return 0, err
	}

	if requestBody.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return 0, errors.New("webhook event is not \"user.upgraded\"")
	}

	return requestBody.Data.UserId, nil
}
