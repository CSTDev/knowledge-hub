package main

import (
	"log"
	"net/http"

	"github.com/cstdev/knowledge-hub/apps/knowledge/knowledge"
	"github.com/rs/cors"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}

func main() {
	db := "mongo:\\\\"
	var service = &knowledge.WebService{DB: &db}

	router := knowledge.NewRouter(service)

	log.Fatal(http.ListenAndServe(":8000", setupGlobalMiddleware(router)))
}
