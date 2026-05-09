package noymw_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ras0q/noy"
	"github.com/ras0q/noy/noymw"
)

func ExampleChain() {
	type User struct {
		Name string
	}
	type State struct {
		CurrentUser *User
	}

	users := map[string]*User{
		"secret": {Name: "Ada"},
	}
	authenticate := func(next noy.HandlerFunc[State]) noy.HandlerFunc[State] {
		return func(state *State, w http.ResponseWriter, r *http.Request) {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			user, ok := users[token]
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			state.CurrentUser = user
			next(state, w, r)
		}
	}
	securityHeaders := noymw.FromStd[State](func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			next.ServeHTTP(w, r)
		})
	})
	handler := noymw.Chain(authenticate, securityHeaders)(func(state *State, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello %s", state.CurrentUser.Name)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	res := httptest.NewRecorder()

	handler(&State{}, res, req)

	fmt.Println(res.Body.String())
	fmt.Println(res.Header().Get("X-Content-Type-Options"))

	// Output:
	// hello Ada
	// nosniff
}
