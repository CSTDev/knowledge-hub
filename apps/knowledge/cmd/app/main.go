package main

import (
	"net/http"

	"github.com/cstdev/knowledge-hub/apps/knowledge/database"
	"github.com/cstdev/knowledge-hub/apps/knowledge/knowledge"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}

func main() {
	dbURL := "172.17.0.2"
	dbName := "knowledge-hub"
	dbCollection := "records"

	log.SetLevel(log.DebugLevel)

	db := &database.MongoDB{
		URL:        dbURL,
		Database:   dbName,
		Collection: dbCollection,
	}

	var service = &knowledge.WebService{DB: db}

	router := knowledge.NewRouter(service)
	log.WithField("port", 8000).Info("Starting server")
	log.Fatal(http.ListenAndServe(":8000", setupGlobalMiddleware(router)))
}
