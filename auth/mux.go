package auth

import (
	"context"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Mux is a request multiplexer
type Mux struct {
	path string

	providers map[string]*Provider
}

// Provider bootstraps oauth2.Config and Provider
type Provider struct {
	cfg      oauth2.Config
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
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

// RegisterProvider attempts to create a new OAuth2 Provider
func (mx *Mux) RegisterProvider(tenantID string) {
	authSvc := viper.GetString("services.auth-service.host") + viper.GetString("services.auth-service.port")
	discoveryURL := authSvc + viper.GetString("services.auth-service.discovery-url") + tenantID
	provider, err := oidc.NewProvider(context.Background(), discoveryURL)
	if err != nil {
		// handle err
	}

	oauth2Cfg := oauth2.Config{
		ClientID:     viper.GetString("oauth.client-id"),
		ClientSecret: viper.GetString("oauth.client-secret"),
		RedirectURL:  viper.GetString("oauth.redirect-url"),
		Endpoint:     provider.Endpoint(),
		Scopes:       append(viper.GetStringSlice("oauth.scopes"), oidc.ScopeOfflineAccess, oidc.ScopeOpenID),
	}

	oidcCfg := &oidc.Config{
		ClientID: viper.GetString("oauth.client-id"),
	}

	tokenVerifier := provider.Verifier(oidcCfg)

	p := &Provider{
		cfg:      oauth2Cfg,
		provider: provider,
		verifier: tokenVerifier,
	}

	mx.providers[tenantID] = p
}

// Handler bootstraps package routes and their respective HTTP HandlerFuncs
// returning a standard http.Handler interface to be served independenly or mounted in a routing chain
func (mx *Mux) Handler() http.Handler {
	handler := chi.NewMux()

	handler.HandleFunc("/authorize", mx.authorize)
	handler.Get("/callback", mx.callback)
	handler.Get("/expired", mx.expired)
	handler.Get("/health", mx.health)
	handler.Post("/login", mx.login)
	handler.Get("/logout", mx.logout)
	handler.Get("/token", mx.token)

	return handler
}
