package main

import (
	"go-api-demo/handlers"
	"go-api-demo/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Main entry point.
func main() {

	// // OLD Beginner way
	// // Always needs the path and a func of signature responseWriter and request as params
	// // Registers the handler function for this specific path
	// http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
	// 	// 	/users handles both GET and POST.
	// 	if r.Method == http.MethodGet {
	// 		getUsers(w, r)
	// 	} else if r.Method == http.MethodPost {
	// 		createUsers(w, r)
	// 	} else {
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

	// // /user handles single-user retrieval.
	// http.HandleFunc("/user", getUserByID)

	// Improved way using Router
	// Now the router handles method + path dispatch.
	r := chi.NewRouter()

	// Global middleware
	// Every request goes through logging and recovery before reaching your handlers.
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	// List users
	r.Get("/users", handlers.GetUsers)

	// Create user
	r.Post("/users", handlers.CreateUsers)

	// Get user by ID
	r.Get("/users/{"+handlers.ParamID+"}", handlers.GetUserByID)
	// 	If a client calls /users/42, chi:
	// - Matches the URL to the route template.
	// - Captures the segment that matches {id} → "42".
	// - Makes it accessible via chi.URLParam(r, "id")

	// start HTTP server, crash if error
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r)) // <- Always remember to pass router here
}
