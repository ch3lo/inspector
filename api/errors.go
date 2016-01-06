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

type SerializationError struct {
	codeAndMessage
	Detail string `json:"detail"`
}

func NewSerializationError(d string) SerializationError {
	return SerializationError{
		codeAndMessage{Code: 400, Message: "Servicio no existe"},
		d,
	}
}

type UnknownError struct {
	codeAndMessage
	detail string
}

func NewUnknownError(d string) UnknownError {
	return UnknownError{
		codeAndMessage{Code: 500, Message: "Error desconocido"},
		d,
	}
}

type ContainerNotFound struct {
	codeAndMessage
}

func NewContainerNotFound() ContainerNotFound {
	return ContainerNotFound{
		codeAndMessage{Code: 404, Message: "El contenedor no existe"},
	}
}
