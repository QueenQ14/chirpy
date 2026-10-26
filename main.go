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
	apiCfg := apiConfig{}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	router.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	router.HandleFunc("GET /metrics", apiCfg.fileserverHitCount)
	router.HandleFunc("GET /healthz", healthCheck)
	router.HandleFunc("POST /reset", apiCfg.hitcountReset)
	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Fatal(server.ListenAndServe())

}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
