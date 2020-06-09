package openid

import (
	"bytes"
	"context"

	"github.com/coreos/go-oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// AuthenticationProvider leverages oauth2.Config and *oidc.Provider to perform  OAuth 2.0
// and OpenID Connect protocol related functions
type AuthenticationProvider struct {
	Config   oauth2.Config
	Provider *oidc.Provider
}

// NewAuthenticationProvider creates a new OAuth2 OpenID Connect Provider
func NewAuthenticationProvider(discoveryURL string) *AuthenticationProvider {
	provider, err := oidc.NewProvider(context.Background(), discoveryURL)
	if err != nil {
		panic(err)
	}

	cfg := oauth2.Config{
		ClientID:     viper.GetString("oauth.client-id"),
		ClientSecret: viper.GetString("oauth.client-secret"),
		RedirectURL:  viper.GetString("oauth.redirect-url"),
		Endpoint:     provider.Endpoint(),
		Scopes:       append(viper.GetStringSlice("oauth.scopes"), oidc.ScopeOfflineAccess, oidc.ScopeOpenID),
	}

	return &AuthenticationProvider{
		Config:   cfg,
		Provider: provider,
	}
}

// Exchange -
func (p *AuthenticationProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.Config.Exchange(ctx, code)
}

// Verify -
func (p *AuthenticationProvider) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	oidcCfg := &oidc.Config{
		ClientID: viper.GetString("oauth.client-id"), // Make clientID a field on receiver
	}

	return p.Provider.Verifier(oidcCfg).Verify(ctx, rawIDToken)
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
