package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adamhu714/chirpy/internal/database"
)

type errorStruct struct {
	Error string `json:"error"`
}

func handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	email, password, err := validateUser(w, r)
	if err != nil {
		return
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerPostUsers - Error while connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = checkEmailUsed(w, email, db)
	if err != nil {
		log.Printf("create user request email already used: %s", err.Error())
		return
	}

	err = db.CreateUser(email, password)
	if err != nil {
		log.Printf("handlerPostUsers - Error while creating user: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		log.Printf("Error retrieving users from database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respUser := users[len(users)-1]
	respUserNoPass := struct {
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		Id:          respUser.Id,
		Email:       respUser.Email,
		IsChirpyRed: respUser.IsChirpyRed,
	}

	data, err := json.Marshal(respUserNoPass)
	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func checkEmailUsed(w http.ResponseWriter, email string, db *database.DB) error {
	users, err := db.GetUsers()
	if err != nil {
		log.Printf("Error retrieving users from database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("error during retrieval of users from database")
	}

	for _, user := range users {
		if user.Email == email {
			respBody := errorStruct{
				Error: "email used",
			}
			respondWithJSON(w, http.StatusBadRequest, respBody)
			return errors.New("email used")
		}
	}
	return nil
}

func validateUser(w http.ResponseWriter, r *http.Request) (string, string, error) {

	type requestParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		return "", "", err
	}

	if len(requestBody.Email) == 0 {
		respBody := errorStruct{
			Error: "Email not provided",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return "", "", errors.New("email message not provided")
	}

	if len(requestBody.Password) == 0 {
		respBody := errorStruct{
			Error: "Password not provided",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return "", "", errors.New("password not provided")
	}

	return requestBody.Email, requestBody.Password, nil
}
