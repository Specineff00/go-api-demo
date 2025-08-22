package middleware

import (
	"net/http"
)

// Logging Middleware
// next http.Handler → the handler being wrapped i.e. getUser, createUser etc
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r) // Call the actual handler
	})
}
