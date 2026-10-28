package main

import (
	"net/http"
	"time"

	"github.com/QueenQ14/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	rows, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	if rows.RevokedAt.Valid == true {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	jwtToken, err := auth.MakeJWT(rows.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating JWT")
		return
	}
	respondWithJson(w, http.StatusOK, response{Token: jwtToken})
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	err = cfg.db.RevokeTokenFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating revoked status for token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
