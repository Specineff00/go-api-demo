package handlers

import (
	"encoding/json"
	"go-api-demo/models"
	"go-api-demo/repositories"
	"go-api-demo/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var userRepo repositories.UserRepository

func InitHandlers(repo repositories.UserRepository) {
	userRepo = repo
}

const ParamID = "id"

// Handlers below
// Handlers are interfaces which functions automatically conform to if you implement the funcs no need to state conformance
// func(w http.ResponseWriter, r *http.Request)

// GET /users → list all users
// http.ResponseWriter like a URLResponse writer (you write the HTTP response to it)
// http.Request pointer to request (like URLRequest).
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userRepo.GetAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

// POST /users → create a new user
func CreateUsers(w http.ResponseWriter, r *http.Request) {
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

	if utils.IsEmpty(user.Name) {
		utils.WriteError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if utils.IsEmpty(user.Email) {
		utils.WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Create a user in DB
	err := userRepo.Create(&user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	// Sends back the response to client
	utils.WriteJSON(w, http.StatusCreated, user)
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

	if utils.IsEmpty(idStr) {
		utils.WriteError(w, http.StatusBadRequest, "Missing ID parameter")
		return
	}

	// Convert to string to int
	id, err := strconv.Atoi(idStr) // can give an error as well, hence two return types
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	user, err := userRepo.GetByID(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "User not found")
		return
	}

	// Send back response (user)
	utils.WriteJSON(w, http.StatusOK, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, ParamID)

	if utils.IsEmpty(idStr) {
		utils.WriteJSON(w, http.StatusBadRequest, "Missing ID Parameters")
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = userRepo.Delete(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, ParamID)

	if utils.IsEmpty(idStr) {
		utils.WriteError(w, http.StatusBadRequest, "Missing ID Parameter")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Parse the request body
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user.ID = id

	// Basic validation
	if utils.IsEmpty(user.Name) {
		utils.WriteError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if utils.IsEmpty(user.Email) {
		utils.WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Try to update the user
	err = userRepo.Update(&user)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}
