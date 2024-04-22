package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adamhu714/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerPostRefresh(w http.ResponseWriter, r *http.Request) {
	id, err := cfg.GetIDFromRefreshToken(w, r)
	if err != nil {
		log.Printf("handlePutUsers - Error getting id from token: %s", err.Error())
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
		log.Printf("Error retrieving users from database: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenAccessString, err := cfg.CreateToken(users[id-1], 60*60, "chirpy-access")
	if err != nil {
		log.Printf("handlerPostLogin - Error creating jwt token: %s", err.Error())
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: tokenAccessString,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (cfg *apiConfig) GetIDFromRefreshToken(w http.ResponseWriter, r *http.Request) (int, error) {
	authHeaderContent := r.Header.Get("Authorization")
	if len(authHeaderContent) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("bad jwt token provided")
	}

	revoked, err := checkIfTokenRevoked(authHeaderContent[7:])
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return 0, errors.New("error checking if token revoked")
	}

	if revoked {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("jwt token provided has been revoked")
	}

	token, err := jwt.ParseWithClaims(
		authHeaderContent[7:],
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, err
	}

	if claims.Issuer != "chirpy-refresh" {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("jwt token provided is not refresh type")
	}

	if claims.ExpiresAt.Unix() < time.Now().UTC().Unix() {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("expired jwt token provided")
	}

	subject := claims.Subject

	id, err := strconv.Atoi(subject)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, err
	}

	return id, nil
}

func checkIfTokenRevoked(token string) (bool, error) {

	db, err := database.NewDB("database.json")
	if err != nil {
		msg := fmt.Sprintf("checkIfTokenRevoked - error while connecting database: %s", err.Error())
		return false, errors.New(msg)
	}

	revokedTokens, err := db.GetRevokedTokens()
	if err != nil {
		msg := fmt.Sprintf("checkIfTokenRevoked - error retrieving revoked tokens from database: %s", err.Error())
		return false, errors.New(msg)
	}

	tokenStatus, ok := revokedTokens[token]
	if !ok {
		return false, nil
	}

	if tokenStatus {
		return true, nil
	}

	return false, nil
}
