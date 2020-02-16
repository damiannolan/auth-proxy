package openid

import (
	"bytes"
	"context"

	"github.com/coreos/go-oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Provider is a simple struct which merges oauth2.Config and *oidc.Provider
type Provider struct {
	*oidc.Provider
	config oauth2.Config
}

// NewProvider creates a new OAuth2 OpenID Connect Provider
func NewProvider(discoveryURL string) *Provider {
	p := new(Provider)
	provider, err := oidc.NewProvider(context.Background(), discoveryURL)
	if err != nil {
		panic(err)
	}

	config := oauth2.Config{
		ClientID:     viper.GetString("oauth.client-id"),
		ClientSecret: viper.GetString("oauth.client-secret"),
		RedirectURL:  viper.GetString("oauth.redirect-url"),
		Endpoint:     p.Endpoint(),
		Scopes:       append(viper.GetStringSlice("oauth.scopes"), oidc.ScopeOfflineAccess, oidc.ScopeOpenID),
	}

	p.Provider = provider
	p.config = config
	return p
}

// BuildDiscoveryURL helper function
func BuildDiscoveryURL(realmID string) string {
	var buf bytes.Buffer
	buf.WriteString(viper.GetString("services.auth-service.host"))
	buf.WriteString(":")
	buf.WriteString(viper.GetString("services.auth-service.port"))
	buf.WriteString(viper.GetString("services.auth-service.discovery-url"))
	buf.WriteString(realmID)

	return buf.String()
}
