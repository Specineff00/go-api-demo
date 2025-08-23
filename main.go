package main

import (
	"go-api-demo/database"
	"go-api-demo/handlers"
	"go-api-demo/middleware"
	"go-api-demo/repositories"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Main entry point.
func main() {

	// Database setup
	database.InitDB()
	database.CreateTables()
	database.SeedData()

	// Init the repo
	userRepo := repositories.NewRespository()

	// Pass in repo to handlers
	handlers.InitHandlers(userRepo)

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
	// Init
	r := chi.NewRouter()

	// Global middleware
	// Every request goes through logging and recovery before reaching your handlers.
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)

	// List users
	r.Get("/users", handlers.GetUsers)
	r.Post("/users", handlers.CreateUsers)
	r.Get("/users/{"+handlers.ParamID+"}", handlers.GetUserByID)
	// 	If a client calls /users/42, chi:
	// - Matches the URL to the route template.
	// - Captures the segment that matches {id} → "42".
	// - Makes it accessible via chi.URLParam(r, "id")

	// start HTTP server, crash if error
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r)) // <- Always remember to pass router here
}
