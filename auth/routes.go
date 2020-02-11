package auth

import (
	"net/http"

	"github.com/damiannolan/auth-proxy/tenant"
)

func (mx *Mux) authorize(w http.ResponseWriter, req *http.Request) {
	tenantID, ok := tenant.FromContext(req.Context())
	if !ok {
		// redirect
	}

	provider := mx.providers[tenantID]
	http.Redirect(w, req, provider.cfg.AuthCodeURL(req.URL.Query().Get("state")), http.StatusTemporaryRedirect)
}

func (mx *Mux) callback(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) expired(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) health(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) login(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) logout(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) token(w http.ResponseWriter, req *http.Request) {

}
