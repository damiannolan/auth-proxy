package middleware

import (
	"net/http"

	"github.com/damiannolan/auth-proxy/tenant"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TenantDiscoveryMiddleware checks for the tenantID in the request cookie and appends it to the request context
// Otherwise redirect for tenant identification
func TenantDiscoveryMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("tenantId")
			if err != nil {
				log.WithError(err).Trace("redirecting for tenant identification")
				redirectURL := viper.GetString("services.tenancy-service.redirect.url")
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			}

			tenantID := cookie.Value
			ctx := tenant.NewContext(r.Context(), tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
