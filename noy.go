package noy

import "net/http"

// ServeMux is a typed wrapper around http.ServeMux.
//
// Handlers registered with HandleFunc receive a new *State for every request.
// The mux still implements http.Handler, so it can be used anywhere a standard
// net/http handler is accepted.
type ServeMux[State any] struct {
	mux *http.ServeMux
}

// NewServeMux creates a ServeMux whose typed handlers receive request-local
// state values of type State.
func NewServeMux[State any]() *ServeMux[State] {
	mux := http.NewServeMux()

	return &ServeMux[State]{
		mux: mux,
	}
}

var _ http.Handler = (*ServeMux[any])(nil)

// ServeHTTP implements http.Handler.
func (m *ServeMux[State]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

// Handle registers a standard http.Handler for pattern.
//
// Use Handle for handlers that do not need typed state, or for mounting another
// mux that already owns its handler behavior.
func (m *ServeMux[State]) Handle(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, handler)
}

// HandleFunc registers a typed handler for pattern.
//
// Each request receives a fresh *State. Middleware and the final handler share
// that pointer for the lifetime of the request.
func (m *ServeMux[State]) HandleFunc(pattern string, handler HandlerFunc[State]) {
	m.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		state := new(State)
		handler(state, w, r)
	})
}

// HandleStdFunc registers a standard http.HandlerFunc for pattern.
//
// Use HandleStdFunc for endpoints that should keep the ordinary net/http
// function signature.
func (m *ServeMux[State]) HandleStdFunc(pattern string, handler http.HandlerFunc) {
	m.mux.HandleFunc(pattern, handler)
}

// HandlerFunc is a typed HTTP handler.
//
// The state pointer is request-local and is passed before the ordinary
// http.ResponseWriter and *http.Request arguments.
type HandlerFunc[State any] func(state *State, w http.ResponseWriter, r *http.Request)
