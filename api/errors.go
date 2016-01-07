package api

import "fmt"

type apiError interface {
	GetCode() int
	GetMessage() string
}

type codeAndMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e codeAndMessage) GetCode() int {
	return e.Code
}

func (e codeAndMessage) GetMessage() string {
	return e.Message
}

func (e codeAndMessage) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Message)
}

// UnknownError error generico que no ha posido ser manejado
type UnknownError struct {
	codeAndMessage
	detail string
}

// NewUnknownError crea un error de tipo UnknownError
func NewUnknownError(d string) UnknownError {
	return UnknownError{
		codeAndMessage{Code: 500, Message: "Error desconocido"},
		d,
	}
}

// ContainerNotFound error generado cuando no existe un contenedor
type ContainerNotFound struct {
	codeAndMessage
}

// NewContainerNotFound Crea un error de tipo ContainerNotFound
func NewContainerNotFound(id string) ContainerNotFound {
	return ContainerNotFound{
		codeAndMessage{Code: 404, Message: fmt.Sprintf("El contenedor %s no existe", id)},
	}
}
