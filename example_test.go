package noy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ras0q/noy"
)

func ExampleServeMux_HandleFunc() {
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
	mux := noy.NewServeMux[State]()
	mux.HandleFunc("GET /me", authenticate(func(state *State, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello %s", state.CurrentUser.Name)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer secret")
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	fmt.Println(res.Body.String())

	// Output:
	// hello Ada
}
