package utilities

import (
	"log"
	"net/http"

	"WebApp/src/middlewares"

	"github.com/google/uuid"
)

// GenerateUUID Utilities function
func GenerateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return id.String()
}

func compose(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, handlerFunc := range middlewares {
		handler = handlerFunc(handler)
	}
	return handler
}

// CallCompose method
func CallCompose(handler http.HandlerFunc) http.HandlerFunc {
	return compose(handler, middlewares.Logging, middlewares.ElapsedTimeForRequest)
}
