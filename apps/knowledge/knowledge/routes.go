package knowledge

import (
	"net/http"

	"github.com/cstdev/knowledge-hub/apps/knowledge/types"
	"github.com/gorilla/mux"
)

// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes []Route

func initRoutes(service types.Service) {
	routes = Routes{
		Route{
			"HealthCheck",
			"GET",
			"/",
			service.HealthCheck(),
		},
		Route{
			"CreateRecord",
			"POST",
			"/record",
			service.NewRecord(),
		}, Route{
			"SearchRecord",
			"GET",
			"/record",
			service.Search(),
		},
	}
}

// NewRouter takes a Service and creates an mux.Router
// Is uses the methods of the Service to associate the handlers
// to their implementations
func NewRouter(s types.Service) *mux.Router {
	initRoutes(s)
	router := mux.NewRouter().StrictSlash(true)

	sub := router.PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		sub.HandleFunc(route.Pattern, route.HandlerFunc).Name(route.Name).Methods(route.Method)
	}
	return router
}
