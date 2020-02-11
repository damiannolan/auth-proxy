package auth

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Mux is a request multiplexer
type Mux struct {
	path string
}

// NewMux returns a new http request multiplexer for the given path
func NewMux(path string) *Mux {
	mux := &Mux{
		path: path,
	}

	return mux
}

// Path returns the configured path
func (mx *Mux) Path() string {
	return mx.path
}

// Handler bootstraps package routes and their respective HTTP HandlerFuncs
// returning a standard http.Handler interface to be served independenly or mounted in a routing chain
func (mx *Mux) Handler() http.Handler {
	handler := chi.NewMux()

	handler.Get("/authorize", mx.authorize)
	handler.Get("/callback", mx.callback)
	handler.Get("/expired", mx.expired)
	handler.Get("/health", mx.health)
	handler.Get("/login", mx.login)
	handler.Get("/logout", mx.logout)
	handler.Get("/token", mx.token)

	return handler
}
