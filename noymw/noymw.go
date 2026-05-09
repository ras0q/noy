package noymw

import (
	"net/http"

	"github.com/ras0q/noy"
)

// FromStd adapts standard net/http middleware to typed noy middleware.
//
// The adapted middleware receives a temporary http.Handler that calls the next
// typed handler with the same request-local state pointer.
func FromStd[State any](httpMiddleware func(http.Handler) http.Handler) func(noy.HandlerFunc[State]) noy.HandlerFunc[State] {
	return func(next noy.HandlerFunc[State]) noy.HandlerFunc[State] {
		return func(state *State, w http.ResponseWriter, r *http.Request) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(state, w, r)
			})
			httpMiddleware(handler).ServeHTTP(w, r)
		}
	}
}

// Chain composes typed middleware in the order provided.
//
// Chain(a, b, c)(handler) runs a, then b, then c, then handler.
func Chain[State any](middlewares ...func(noy.HandlerFunc[State]) noy.HandlerFunc[State]) func(noy.HandlerFunc[State]) noy.HandlerFunc[State] {
	return func(next noy.HandlerFunc[State]) noy.HandlerFunc[State] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
