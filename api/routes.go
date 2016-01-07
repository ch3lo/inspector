package api

import (
	"github.com/ch3lo/inspector/logger"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
)

type appContext struct {
	hostIP string
	client *docker.Client
}

var routesMap = map[string]map[string]serviceHandler{
	"GET": {
		"/container/{id}": getInspectContainer,
	},
}

func routes(config Configuration, sts *stats.Stats) *mux.Router {
	ctx := &appContext{
		hostIP: config.Advertise,
	}

	logger.Instance().Debugf("Configurando API de Docker con los parametros %+v", config)

	var err error

	if config.TLSVerify {
		ctx.client, err = docker.NewTLSClient(config.Address, config.TLSCert, config.TLSKey, config.TLSCacert)
	} else {
		ctx.client, err = docker.NewClient(config.Address)
	}
	if err != nil {
		logger.Instance().Fatalln("Error al crear el cliente")
	}

	router := mux.NewRouter()

	router.Handle("/stats", &statsHandler{sts}).Methods("GET")

	// API v1
	v1Services := router.PathPrefix("/api/v1").Subrouter()

	for method, mappings := range routesMap {
		for path, h := range mappings {
			v1Services.Handle(path, errorHandler{h, ctx}).Methods(method)
		}
	}

	return router
}
