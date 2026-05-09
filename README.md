# noy

[![Go Reference](https://pkg.go.dev/badge/github.com/ras0q/noy.svg)](https://pkg.go.dev/github.com/ras0q/noy)

`noy` is a tiny typed layer over Go's standard [`net/http`](https://pkg.go.dev/net/http) routing.

It keeps the standard library routing model, then adds one request-local state pointer that middleware and handlers can share without `context.WithValue`.

## Install

````sh
go get github.com/ras0q/noy
````

## What noy changes

In plain `net/http`, request-local application data usually moves through `context.Context`.

````go
// In authentication middleware
ctx := context.WithValue(currentUserKey, authenticatedUser)
r = r.WithContext(ctx)

// In each handler
user := r.Context().Value(currentUserKey).(*User)
````

That works, but the compiler cannot prove that the key exists or that the value has the expected type.

With `noy`, your application declares one state type for request-local data.

````go
type State struct {
	CurrentUser *User
	RequestID   string
}
````

Every typed middleware and handler receives `*State`.

See the [`_examples/`](./_examples/) directory for more examples.

## Request flow

For an authenticated API endpoint, the flow looks like this:

1. A request enters `noy.NewServeMux[State]()`.
2. `noy` creates a fresh `*State` for that request.
3. Auth middleware validates the `Authorization` header and sets `state.CurrentUser`.
4. Other middleware can read or add request-local data, such as `state.RequestID`.
5. The endpoint handler reads `state.CurrentUser` directly and writes the response.

The handler does not need to know how authentication loaded the user, and the authentication middleware does not need to hide values behind context keys.

````go
func profile(state *State, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", state.CurrentUser.Name)
}
````

## Middleware

Typed middleware wraps `noy.HandlerFunc[State]`, so middleware can populate state before the final handler runs.

`noymw` provides small helpers for chaining typed middleware and adapting ordinary `func(http.Handler) http.Handler` middleware. See the [`noymw` package documentation](https://pkg.go.dev/github.com/ras0q/noy/noymw) for details.

## Standard net/http interop

`noy` stays close to the standard library:

- route patterns and path values are handled by `http.ServeMux`
- the mux can be passed anywhere an `http.Handler` is accepted
- ordinary `http.Handler` and `http.HandlerFunc` values can still be mounted
- existing standard middleware can be adapted with `noymw`
