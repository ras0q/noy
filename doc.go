// Package noy adds typed, request-local state to applications built on net/http.
//
// The core idea is that an application defines one state type for values that
// live for a single request. Middleware can populate that state with data such
// as authenticated users, request IDs, or loaded domain objects, and handlers
// can read those values through ordinary Go fields.
//
// noy keeps routing and handler composition close to the standard library. The
// additional convention is explicit state passing, which keeps request-local
// data visible to the compiler instead of hiding it behind context keys.
package noy
