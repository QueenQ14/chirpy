package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/QueenQ14/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set in the environment")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set in the environment")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	const filepathRoot = "."
	const port = "8080"
	router := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
		jwtSecret:      jwtSecret,
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	router.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	router.HandleFunc("GET /admin/metrics", apiCfg.fileserverHitCount)
	router.HandleFunc("GET /api/healthz", healthCheck)
	router.HandleFunc("POST /admin/reset", apiCfg.hitcountReset)
	router.HandleFunc("POST /api/chirps", apiCfg.postChirp)
	router.HandleFunc("GET /api/chirps", apiCfg.listChirps)
	router.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByID)
	router.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	router.HandleFunc("POST /api/users", apiCfg.createUser)
	router.HandleFunc("PUT /api/users", apiCfg.updateUser)
	router.HandleFunc("POST /api/login", apiCfg.handleLogin)
	router.HandleFunc("POST /api/refresh", apiCfg.refreshToken)
	router.HandleFunc("POST /api/revoke", apiCfg.revokeToken)
	router.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUser)
	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	fmt.Println("The server has started...")
	log.Fatal(server.ListenAndServe())

}
