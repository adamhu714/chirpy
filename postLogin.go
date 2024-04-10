package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adamhu714/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerPostLogin(w http.ResponseWriter, r *http.Request) {
	email, password, expiresInSeconds, err := validateUserLogin(w, r)
	if err != nil {
		return
	}

	if expiresInSeconds == 0 {
		expiresInSeconds = 24 * 60 * 60
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Printf("handlerPostLogin - Error while connecting database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := findUser(w, email, password, db)
	if err != nil {
		log.Printf("handlerPostLogin - Error finding user: %s", err.Error())
		return
	}

	// Add funcctin here

	respUserNoPass := struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    user.Id,
		Email: user.Email,
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

func validateUserLogin(w http.ResponseWriter, r *http.Request) (string, string, int, error) {

	type requestParams struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
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
		return "", "", 0, err
	}

	if len(requestBody.Email) == 0 {
		respBody := errorStruct{
			Error: "Email not provided",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return "", "", 0, errors.New("email message not provided")
	}

	if len(requestBody.Password) == 0 {
		respBody := errorStruct{
			Error: "Password not provided",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return "", "", 0, errors.New("password not provided")
	}

	return requestBody.Email, requestBody.Password, requestBody.Expires_in_seconds, nil
}
