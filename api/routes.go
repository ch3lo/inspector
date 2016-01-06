package api

import (
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
)

var routesMap = map[string]map[string]serviceHandler{
	"GET": {
		"/container/{id}": getInspectContainer,
	},
}

func routes(sts *stats.Stats) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/stats", &statsHandler{sts}).Methods("GET")

	// API v1
	v1Services := router.PathPrefix("/api/v1").Subrouter()

	for method, mappings := range routesMap {
		for path, h := range mappings {
			v1Services.Handle(path, errorHandler(h)).Methods(method)
		}
	}

	return router
}
