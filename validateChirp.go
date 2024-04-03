package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type requestParams struct {
		Body string `json:"body"`
	}
	type errorStruct struct {
		Error string `json:"error"`
	}
	type cleanParams struct {
		CleanedBody string `json:"cleaned_body"`
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
		return
	}

	if len(requestBody.Body) > 140 {
		respBody := errorStruct{
			Error: "Chirp is too long",
		}
		respondWithJSON(w, http.StatusBadRequest, respBody)
		return
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

	respBody := cleanParams{
		CleanedBody: cleanedBody,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}

func respondWithJSON(w http.ResponseWriter, code int, respBody interface{}) {
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
	w.WriteHeader(code)
}
