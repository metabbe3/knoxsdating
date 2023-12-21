package main

import (
	"log"
	"net/http"

	"github.com/metabbe3/knoxsdating/pkg/routes"
)

func main() {
	// Initialize all routes
	router := routes.InitializeRoutes()

	// Add logging middleware to log requests
	router.Use(loggingMiddleware)

	// Start the HTTP server
	port := ":8080" // Default port
	log.Printf("Server listening on %s\n", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

// loggingMiddleware is a middleware function to log requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
