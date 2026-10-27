package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Body string `json:"cleaned_body"`
	}
	type request struct {
		Body string `json:"body"`
	}
	profaneWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "there was an error decoding the request")
		return
	}
	if len(req.Body) <= 140 {
		var filtered_body []string
		for _, word := range strings.Fields(req.Body) {
			if _, ok := profaneWords[strings.ToLower(word)]; !ok {
				filtered_body = append(filtered_body, word)
			} else {
				filtered_body = append(filtered_body, "****")
			}
		}
		resBody := strings.Join(filtered_body, " ")
		respondWithJson(w, http.StatusOK, response{Body: resBody})
	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}

}
