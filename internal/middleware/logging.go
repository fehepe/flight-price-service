package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging logs request method, URI and duration.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s completed in %v", r.Method, r.RequestURI, time.Since(start))
	})
}
