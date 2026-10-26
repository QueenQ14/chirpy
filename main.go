package main

import (
	"net/http"
)

func main() {
	router := http.NewServeMux()
	server := http.Server{
		Handler: router,
		Addr:    ":8080",
	}
	server.ListenAndServe()

}
