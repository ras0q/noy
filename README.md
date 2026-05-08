# noy

A simple, lightweight wrapper for Go's standard [`http.ServeMux`](https://pkg.go.dev/net/http#ServeMux) that brings strongly-typed, request-local state to your HTTP handlers.

## Why use `noy`?

- Use Go Generics for strictly typed request state instead of `context.WithValue`.
- Pass data easily by updating the `State` struct directly.
- Integrate seamlessly with the standard `net/http` library.

## Example

See the [`_examples/`](./_examples) directory for more examples.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/ras0q/noy"
)

type State struct {
	UserID int
}

func main() {
	mux := noy.NewServeMux[State]()

	mux.HandleFunc("/", func(state *State, w http.ResponseWriter, r *http.Request) {
		state.UserID = 42 // Access and mutate state directly
		fmt.Fprintf(w, "Hello, User %d!", state.UserID)
	})

	http.ListenAndServe(":8080", mux)
}
```
