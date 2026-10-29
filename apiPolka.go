package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err := cfg.db.UpdateUserStatusByID(r.Context(), req.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)
}
