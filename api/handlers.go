package api

import (
	"fmt"
	"net/http"

	"github.com/ch3lo/inspector/api/types"
	"github.com/ch3lo/inspector/util"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

type serviceHandler func(w http.ResponseWriter, r *http.Request) error

type errorHandler serviceHandler

func (eh errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := eh(w, r); err != nil {
		util.Log.Errorln(err)
		if err2, ok := err.(apiError); ok {
			jsonRenderer(w, err2)
			return
		}
		jsonRenderer(w, NewUnknownError(err.Error()))
	}
}

type statsHandler struct {
	*stats.Stats
}

func (sh *statsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jsonRenderer(w, sh.Data())
}

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func jsonRenderer(w http.ResponseWriter, i interface{}) {
	rend := render.New()
	rend.JSON(w, http.StatusOK, i)
}

func getInspectContainer(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return NewContainerNotFound()
	}

	container, err := client.InspectContainer(id)
	if err != nil {
		return fmt.Errorf("Problema al obtener el contenedor", err)
	}

	util.Log.Debugf("Container Dump %+v", container)

	ports := make(map[string]types.Port)

	for k, publicMap := range container.NetworkSettings.Ports {
		if len(publicMap) == 0 {
			continue
		}

		p, ok := ports[k.Port()]
		if !ok {
			p = types.Port{
				Type: k.Proto(),
			}
		}
		for _, publicPorts := range publicMap {
			p.Publics = append(p.Publics, publicPorts.HostPort)
		}
		ports[k.Port()] = p
	}

	resp := types.Container{
		ID:     container.ID,
		HostIP: advertise,
		Ports:  ports,
	}

	jsonRenderer(w, map[string]interface{}{
		"status":    http.StatusOK,
		"container": resp})
	return nil
}
