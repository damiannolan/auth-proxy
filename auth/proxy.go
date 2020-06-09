package auth

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/damiannolan/auth-proxy/openid"
	"github.com/damiannolan/auth-proxy/realm"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ProxyService -
type ProxyService struct {
	path      string
	providers map[string]*openid.AuthenticationProvider
	upstream  *httputil.ReverseProxy
}

// NewProxyService -
func NewProxyService() (*ProxyService, error) { // Maybe pass viper instance or config with fields - would be large config
	rawURL := fmt.Sprintf("%s:%s", viper.GetString("services.upstream.host"), viper.GetString("services.upstream.port"))
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	proxy := &ProxyService{
		path:      viper.GetString("oauth.url"),
		providers: make(map[string]*openid.AuthenticationProvider),
		upstream:  httputil.NewSingleHostReverseProxy(url),
	}

	return proxy, nil
}

// Handler bootstraps package routes and their respective HTTP HandlerFuncs
// returning a standard http.Handler interface to be served independenly or mounted in a routing chain
func (svc *ProxyService) Handler() http.Handler {
	mux := chi.NewMux()

	mux.Use(middleware.DefaultLogger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)
	mux.Use(realm.DiscoveryMiddleware())

	// Change the below to be wrapped in a middleware func - mux.Use()
	// mux.HandleFunc("/*", svc.upstream.ServeHTTP)

	mux.Route(svc.Path(), func(r chi.Router) {
		r.HandleFunc("/authorize", svc.authorize)
		r.Get("/callback", svc.callback)
		r.Get("/expired", svc.expired)
		r.Get("/health", svc.health)
		r.Post("/login", svc.login)
		r.Get("/logout", svc.logout)
		r.Get("/token", svc.token)
	})

	return mux
}

// ListenAndServe function
func (svc *ProxyService) ListenAndServe() error {
	host := viper.GetString("services.auth-proxy.host")
	port := viper.GetString(("services.auth-proxy.port"))
	log.WithField("server", host).Info("listening on port ", port)

	return http.ListenAndServe(port, svc.Handler())
}

// Path returns the configured path
func (svc *ProxyService) Path() string {
	return svc.path
}

// RegisterProvider inserts the OpenID Connect Provider into the local cache
func (svc *ProxyService) RegisterProvider(key string, p *openid.AuthenticationProvider) {
	svc.providers[key] = p
}
