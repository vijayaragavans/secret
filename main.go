package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vijayaragavans/secret/api"
)

var Lang = "en"

func withTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a context with a 10-second timeout
		ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
		defer cancel()

		// Replace the request's context with the new one
		r = r.WithContext(ctx)

		// Handle the request
		done := make(chan struct{})
		go func() {
			next.ServeHTTP(w, r)
			close(done)
		}()

		select {
		case <-ctx.Done():
			// If the context is done, return a timeout response
			if ctx.Err() == context.DeadlineExceeded {
			}
		case <-done:
			// Request completed successfully
		}
	})
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/generate", api.Generate).Methods("POST")
	r.HandleFunc("/read/{key}", api.Read).Methods("GET")

	r.Use(withTimeout)

	c := cors.AllowAll()
	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":8080", handler))

}
