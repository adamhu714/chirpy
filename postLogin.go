package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adamhu714/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func handlerPostLogin(w http.ResponseWriter, r *http.Request) {
	email, password, err := validateUser(w, r)
	if err != nil {
		return
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerPostLogin - Error while connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respUser, err := findUser(w, email, password, db)
	if err != nil {
		log.Printf("handlerPostLogin - Error finding user: %s", err.Error())
		return
	}

	respUserNoPass := struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    respUser.Id,
		Email: respUser.Email,
	}
	data, err := json.Marshal(respUserNoPass)
	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func findUser(w http.ResponseWriter, email string, password string, db *database.DB) (database.User, error) {
	users, err := db.GetUsers()
	if err != nil {
		log.Printf("Error retrieving users from database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return database.User{}, errors.New("error during retrieval of users from database")
	}

	for _, user := range users {
		if user.Email == email {
			err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
			if err == nil {
				return user, nil
			}
			respBody := errorStruct{
				Error: "email or password is incorrect",
			}
			respondWithJSON(w, http.StatusUnauthorized, respBody)
			return database.User{}, errors.New("login: password incorrect")
		}
	}

	respBody := errorStruct{
		Error: "email or password is incorrect",
	}
	respondWithJSON(w, http.StatusUnauthorized, respBody)
	return database.User{}, errors.New("login: user not found")
}
