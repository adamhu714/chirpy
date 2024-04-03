package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type requestParams struct {
		Body string `json:"body"`
	}

	type errorStruct struct {
		Error string `json:"error"`
	}

	type validParams struct {
		Valid bool `json:"valid"`
	}

	var requestBody requestParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)

	if err != nil {
		log.Printf("Error while json decoding: %s", err.Error())

		respBody := errorStruct{
			Error: "Something went wrong",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error while json marshalling: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(requestBody.Body) > 140 {
		respBody := errorStruct{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error while json marshalling: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respBody := validParams{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)

	if err != nil {
		log.Printf("Error while json marshalling: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(dat)
	w.WriteHeader(http.StatusOK)
}
