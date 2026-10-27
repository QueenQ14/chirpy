package main

import (
	"encoding/json"
	"net/http"
)

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type jsonValid struct {
		Valid bool `json:"valid"`
	}
	type request struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "there was an error decoding the request")
		return
	}
	if len(req.Body) <= 140 {
		respondWithJson(w, http.StatusOK, jsonValid{Valid: true})
	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}

}
