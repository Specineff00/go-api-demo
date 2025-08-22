package handlers

import (
	"encoding/json"
	"go-api-demo/models"
	"go-api-demo/utils"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

// global state (in-memory DB)
// For learning purposes and not actual DB
var (
	users   = []models.User{} // Slice not an Array (Mutable reference type or dynamic array)
	nextID  = 1               // This is mutable!
	userMUX sync.Mutex        // Prevents Data Races
)

const ParamID = "id"

// Handlers below
// Handlers are interfaces which functions automatically conform to if you implement the funcs no need to state conformance
// func(w http.ResponseWriter, r *http.Request)

// GET /users → list all users
// http.ResponseWriter like a URLResponse writer (you write the HTTP response to it)
// http.Request pointer to request (like URLRequest).
func GetUsers(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, users)
}

// POST /users → create a new user
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	// Anonymous struct → no need to define a type globally.
	var user models.User

	// Decodes the body and stores it in pointer
	// return value is always an optional error (maybe is a pointer?)
	// Error checking
	// If decoding fails, send back a 400 Bad Request.
	// This keeps the err within the scope of the condition. you could separate it but having err
	// expose outside doesnt make for good practice
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if user.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Name is required")
	}

	// lock before mutating global state (like queue.sync).
	userMUX.Lock()
	users = append(users, user)
	nextID++
	userMUX.Unlock()

	utils.WriteJSON(w, http.StatusOK, user)
}

// TODO: why cant we extend string???
// func (s string) IsEmpty() bool {
// 	s == ""
// }

// Get /user/{id} -> get a single user
func GetUserByID(w http.ResponseWriter, r *http.Request) {

	// Why are path params better than queries:
	// More RESTful: /users/42 clearly identifies a resource.
	// Cleaner: no ?id=42 in the URL.
	// Easier to handle nested resources later (/users/42/posts/3).
	idStr := chi.URLParam(r, ParamID)

	// Extracts id query param (?id=123).
	// ids := r.URL.Query().Get("id")

	if idStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing ID parameter")
		return
	}

	// Convert to string to int
	id, err := strconv.Atoi(idStr) // can give an error as well, hence two return types

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	for _, u := range users {
		// If ID matches, return that user.
		if u.ID == id {
			utils.WriteJSON(w, http.StatusOK, u)
			return
		}
	}

	// If loop finishes without finding → 404 Not Found.
	http.NotFound(w, r)
}
