package main

import (
	"net/http"

	"github.com/cstdev/knowledge-hub/apps/knowledge/knowledge"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}

func main() {
	var service = &knowledge.WebService{}

	router := knowledge.NewRouter(service)
	log.Info("starting server on port %d", 8000)
	log.Fatal(http.ListenAndServe(":8000", setupGlobalMiddleware(router)))
}
