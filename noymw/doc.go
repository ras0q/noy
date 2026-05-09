// Package noymw provides middleware helpers for applications using noy.
//
// The package focuses on composing typed middleware and bridging ordinary
// net/http middleware into a pipeline that preserves the same request-local
// state pointer from one layer to the next.
package noymw
