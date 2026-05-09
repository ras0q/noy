package noymw

import (
	"net/http"

	"github.com/ras0q/noy"
)

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

func Chain[State any](middlewares ...func(noy.HandlerFunc[State]) noy.HandlerFunc[State]) func(noy.HandlerFunc[State]) noy.HandlerFunc[State] {
	return func(next noy.HandlerFunc[State]) noy.HandlerFunc[State] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
