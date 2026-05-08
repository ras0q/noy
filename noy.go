package noy

import "net/http"

type ServeMux[State any] struct {
	mux *http.ServeMux
}

func NewServeMux[State any]() *ServeMux[State] {
	mux := http.NewServeMux()

	return &ServeMux[State]{
		mux: mux,
	}
}

var _ http.Handler = (*ServeMux[any])(nil)

func (m *ServeMux[State]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *ServeMux[State]) Handle(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, handler)
}

func (m *ServeMux[State]) HandleFunc(pattern string, handler HandlerFunc[State]) {
	m.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		state := new(State)
		handler(state, w, r)
	})
}

type HandlerFunc[State any] func(state *State, w http.ResponseWriter, r *http.Request)
