package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
)

// Handler holds application state.
type Handler struct {
	// Repositories, services, etc.
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int) // ❌ This is not type-safe
	fmt.Fprintf(w, "Hello, %d!", userID)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	username := r.PathValue("username")
	fmt.Fprintf(w, "User profile for %s (UserID: %d)", username, userID)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := rand.Int()
		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))
		next(w, r)
	}
}

func main() {
	mux := http.NewServeMux()

	h := &Handler{}

	mux.HandleFunc("/", authMiddleware(h.Index))
	{
		usersMux := http.NewServeMux()
		usersMux.HandleFunc("/users/{username}", authMiddleware(h.GetUser))

		mux.Handle("/users/", usersMux)
	}

	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
