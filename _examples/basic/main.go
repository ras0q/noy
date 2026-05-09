package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/ras0q/noy"
	"github.com/ras0q/noy/noymw"
)

// State holds request-local state.
type State struct {
	UserID int
}

// Handler holds application state.
type Handler struct {
	// Repositories, services, etc.
}

func (h *Handler) Index(state *State, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %d!", state.UserID)
}

func (h *Handler) GetUser(state *State, w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	fmt.Fprintf(w, "User profile for %s (UserID: %d)", username, state.UserID)
}

func authMiddleware(next noy.HandlerFunc[State]) noy.HandlerFunc[State] {
	return func(state *State, w http.ResponseWriter, r *http.Request) {
		state.UserID = rand.Int()
		next(state, w, r)
	}
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := noy.NewServeMux[State]()
	h := &Handler{}
	middlewares := noymw.Chain(
		authMiddleware,
		noymw.FromStd[State](loggerMiddleware),
	)

	mux.HandleFunc("/", middlewares(h.Index))
	{
		usersMux := noy.NewServeMux[State]()
		usersMux.HandleFunc("/users/{username}", middlewares(h.GetUser))

		mux.Handle("/users/", usersMux)
	}

	slog.Info("Server starting", "address", ":8080")
	http.ListenAndServe(":8080", mux)
}
