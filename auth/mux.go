package auth

import (
	"net/http"

	"github.com/damiannolan/auth-proxy/openid"
	"github.com/damiannolan/auth-proxy/realm"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Mux is a request multiplexer
type Mux struct {
	path string

	providers map[string]*openid.Provider
}

// NewMux returns a new http request multiplexer for the given path
func NewMux(path string) *Mux {
	mux := &Mux{
		path:      path,
		providers: make(map[string]*openid.Provider),
	}

	return mux
}

// Path returns the configured path
func (mx *Mux) Path() string {
	return mx.path
}

// RegisterProvider inserts the OpenID Connect Provider into the local cache
func (mx *Mux) RegisterProvider(key string, p *openid.Provider) {
	mx.providers[key] = p
}

// Handler bootstraps package routes and their respective HTTP HandlerFuncs
// returning a standard http.Handler interface to be served independenly or mounted in a routing chain
func (mx *Mux) Handler() http.Handler {
	handler := chi.NewMux()

	handler.Use(middleware.DefaultLogger)
	handler.Use(middleware.Recoverer)
	handler.Use(middleware.RequestID)
	handler.Use(realm.DiscoveryMiddleware())

	handler.Route(mx.Path(), func(r chi.Router) {
		r.HandleFunc("/authorize", mx.authorize)
		r.Get("/callback", mx.callback)
		r.Get("/expired", mx.expired)
		r.Get("/health", mx.health)
		r.Post("/login", mx.login)
		r.Get("/logout", mx.logout)
		r.Get("/token", mx.token)
	})

	return handler
}
