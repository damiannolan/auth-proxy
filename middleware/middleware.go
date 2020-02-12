package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/damiannolan/auth-proxy/realm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RealmDiscoveryMiddleware checks for the realmID in the request cookie and appends it to the request context
// Otherwise redirect for realm identification
func RealmDiscoveryMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("realmId")
			if err != nil {
				log.WithError(err).Trace("redirecting for realm identification")

				state := base64.StdEncoding.EncodeToString([]byte(r.URL.RequestURI()))
				redirectURL := fmt.Sprintf("%s?state=%s", viper.GetString("services.tenancy-service.redirect-url"), state)
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			}

			realmID := cookie.Value
			ctx := realm.NewContext(r.Context(), realmID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
