package main

import (
	"net/http"
	"os"

	"github.com/cstdev/knowledge-hub/apps/knowledge/database"
	"github.com/cstdev/knowledge-hub/apps/knowledge/knowledge"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler
	return handleCORS(handler)
}

func main() {
	//dbURL := "172.17.0.2"
	dbName := "knowledge-hub"
	dbCollection := "records"
	fieldCollection := "fields"

	switch logLevel := os.Getenv("LOG_LEVEL"); logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.JSONFormatter{})

	dbURL := os.Getenv("MONGODB_URI")

	if dbURL == "" {
		log.Fatal("$MONGODB_URI must be set")
	}

	db := &database.MongoDB{
		URL:             dbURL,
		Database:        dbName,
		Collection:      dbCollection,
		FieldCollection: fieldCollection,
	}

	var service = &knowledge.WebService{DB: db}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$Port must be set")
	}

	router := knowledge.NewRouter(service)
	log.WithField("port", port).Info("Starting server")
	log.Fatal(http.ListenAndServe(":"+port, setupGlobalMiddleware(router)))
}
