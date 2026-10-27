package main

import (
	"database/sql"
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
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	router.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	router.HandleFunc("GET /admin/metrics", apiCfg.fileserverHitCount)
	router.HandleFunc("GET /api/healthz", healthCheck)
	router.HandleFunc("POST /admin/reset", apiCfg.hitcountReset)
	router.HandleFunc("POST /api/chirps", apiCfg.handleChirp)
	router.HandleFunc("POST /api/users", apiCfg.createUser)
	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Fatal(server.ListenAndServe())

}
