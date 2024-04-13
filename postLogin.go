package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adamhu714/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
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

	tokenString, err := cfg.CreateToken(user, expiresInSeconds)
	if err != nil {
		log.Printf("handlerPostLogin - Error creating jwt token: %s", err.Error())
		return
	}

	respUserNoPass := struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `jaon:"token"`
	}{
		Id:    user.Id,
		Email: user.Email,
		Token: tokenString,
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

func (cfg *apiConfig) CreateToken(user database.User, expiresInSeconds int) (string, error) {
	if expiresInSeconds == 0 || expiresInSeconds > 24*60*60 {
		expiresInSeconds = 24 * 60 * 60
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiresInSeconds))),
		Subject:   strconv.Itoa(user.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
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
