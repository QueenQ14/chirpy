package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	router := http.NewServeMux()
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	router.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	router.HandleFunc("GET /admin/metrics", apiCfg.fileserverHitCount)
	router.HandleFunc("GET /api/healthz", healthCheck)
	router.HandleFunc("POST /admin/reset", apiCfg.hitcountReset)
	router.HandleFunc("/api/validate_chirp", validateChirp)
	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Fatal(server.ListenAndServe())

}
