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
)

func (cfg *apiConfig) handlerPutUsers(w http.ResponseWriter, r *http.Request) {
	id, err := cfg.GetIDFromToken(w, r)
	if err != nil {
		log.Printf("handlePutUsers - Error getting id from token: %s", err.Error())
		return
	}

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

	err = checkEmailUsedPutUsers(w, email, id, db)
	if err != nil {
		log.Printf("create user request email already used: %s", err.Error())
		return
	}

	err = db.UpdateUser(email, password, id)
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

	respUser := users[id]
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
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func checkEmailUsedPutUsers(w http.ResponseWriter, email string, id int, db *database.DB) error {
	users, err := db.GetUsers()
	if err != nil {
		log.Printf("Error retrieving users from database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("error during retrieval of users from database")
	}

	for _, user := range users {
		if user.Id == id {
			if user.Email == email {
				break
			}
			continue
		}
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

func (cfg *apiConfig) GetIDFromToken(w http.ResponseWriter, r *http.Request) (int, error) {
	authHeaderContent := r.Header.Get("Authorization")
	if len(authHeaderContent) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("bad jwt token provided")
	}

	token, err := jwt.ParseWithClaims(
		authHeaderContent[7:],
		jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("bad jwt token provided")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("bad jwt token provided")
	}

	if claims.ExpiresAt.Unix() < time.Now().UTC().Unix() {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("expired jwt token provided")
	}

	subject := claims.Subject

	id, err := strconv.Atoi(subject)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("bad jwt token provided")
	}

	return id, nil
}
