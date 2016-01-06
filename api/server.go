package api

import (
	"github.com/ch3lo/inspector/util"
	"github.com/codegangsta/negroni"
	"github.com/fsouza/go-dockerclient"
	"github.com/rs/cors"
	"github.com/thoas/stats"
)

type Configuration struct {
	Advertise string
	Address   string
	TlsVerify bool
	TlsCacert string
	TlsCert   string
	TlsKey    string
}

var client *docker.Client
var advertise string

func Server(config Configuration) {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST, GET, OPTIONS, PUT, DELETE, UPDATE"},
		AllowedHeaders:   []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		MaxAge:           50,
		AllowCredentials: true,
	})

	advertise = config.Advertise

	var err error
	util.Log.Debugf("Configurando API de Docker con los parametros %+v", config)
	if config.TlsVerify {
		client, err = docker.NewTLSClient(config.Address, config.TlsCert, config.TlsKey, config.TlsCacert)
	} else {
		client, err = docker.NewClient(config.Address)
	}
	if err != nil {
		panic("Error al crear el cliente")
	}

	statsMiddleware := stats.New()

	router := routes(statsMiddleware)

	n := negroni.Classic()
	n.Use(corsMiddleware)
	n.Use(statsMiddleware)
	n.UseHandler(router)

	n.Run(":8080")
}
