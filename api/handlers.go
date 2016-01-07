package api

import (
	"fmt"
	"net/http"

	"github.com/ch3lo/inspector/api/types"
	"github.com/ch3lo/inspector/logger"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
	"github.com/unrolled/render"
)

type serviceHandler func(c *appContext, w http.ResponseWriter, r *http.Request) error

type errorHandler struct {
	handler serviceHandler
	appCtx  *appContext
}

func (eh errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := eh.handler(eh.appCtx, w, r); err != nil {
		logger.Instance().Errorln(err)
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

func jsonRenderer(w http.ResponseWriter, i interface{}) {
	rend := render.New()
	rend.JSON(w, http.StatusOK, i)
}

func getInspectContainer(c *appContext, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return NewContainerNotFound(id)
	}

	container, err := c.client.InspectContainer(id)
	if err != nil {
		switch err.(type) {
		case *docker.NoSuchContainer:
			return NewContainerNotFound(id)
		default:
			return fmt.Errorf("Problema al obtener el contenedor %s", err)
		}

	}

	logger.Instance().Debugf("Container Dump %+v", container)

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
		HostIP: c.hostIP,
		Ports:  ports,
	}

	jsonRenderer(w, map[string]interface{}{
		"status":    http.StatusOK,
		"container": resp})
	return nil
}
