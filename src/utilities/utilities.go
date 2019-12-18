package utilities

import (
	"log"
	"middlewares"
	"net/http"

	"github.com/google/uuid"
)

// Utilities function
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

func CallCompose(handler http.HandlerFunc) http.HandlerFunc {
	return compose(handler, middlewares.Logging, middlewares.ElapsedTimeForRequest)
}
