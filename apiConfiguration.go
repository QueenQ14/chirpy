package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/QueenQ14/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitcountReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "not allowed in non-dev environment")
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(r.Context())
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) fileserverHitCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	data := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(data))
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	req := request{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "there was an error decoding the request")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     req.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}
	res := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJson(w, 201, res)
}
