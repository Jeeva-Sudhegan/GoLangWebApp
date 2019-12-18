package middlewares

import (
	"log"
	"net/http"
	"time"
)

// Middlewares
func Logging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logging: " + r.URL.Path)
		handler(w, r)
	}
}

func ElapsedTimeForRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Println("Elapsed time for the request "+r.URL.Path+" is", time.Since(start))
		}()
		handler(w, r)
	}
}
