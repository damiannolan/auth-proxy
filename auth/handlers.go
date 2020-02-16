package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/damiannolan/auth-proxy/openid"
	"github.com/damiannolan/auth-proxy/realm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (mx *Mux) authorize(w http.ResponseWriter, req *http.Request) {
	realmID, ok := realm.FromContext(req.Context())
	if !ok {
		redirectURL := "getRedirectURL()"
		http.Redirect(w, req, "", http.StatusTemporaryRedirect)
	}

	if provider, ok := mx.providers[realmID]; ok { // "default"
		http.Redirect(w, req, provider.config.AuthCodeURL(req.URL.Query().Get("state")), http.StatusTemporaryRedirect)
	} else {
		// TODO: Decide best approach
		discoveryURL := openid.BuildDiscoveryURL(realmID)
		var buf bytes.Buffer
		buf.WriteString(viper.GetString("services.auth-service.host"))
		buf.WriteString(":")
		buf.WriteString(viper.GetString("services.auth-service.port"))
		buf.WriteString(viper.GetString("services.auth-service.discovery-url"))
		buf.WriteString(realmID)

		provider := openid.NewProvider(discoveryURL)
		mx.providers[realmID] = provider
		http.Redirect(w, req, provider.config.AuthCodeURL(req.URL.Query().Get("state")), http.StatusTemporaryRedirect)
	}
}

func (mx *Mux) callback(w http.ResponseWriter, req *http.Request) {
	realmID, ok := realm.FromContext(req.Context())
	if !ok {
		http.Redirect(w, req, "", http.StatusTemporaryRedirect)
	}

	code := req.URL.Query().Get("code")
	p := mx.providers[realmID]

	token, err := p.config.Exchange(req.Context(), code)
	if err != nil {
		http.Redirect(w, req, "", http.StatusForbidden)
	}

	log.WithField("token", token).Debug("oauth2 token")

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Redirect(w, req, "", http.StatusForbidden)
	}

	idToken, err := p.verifier.Verify(req.Context(), rawIDToken)
	if err != nil {
		http.Redirect(w, req, "", http.StatusForbidden)
	}

	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Redirect(w, req, "", http.StatusInternalServerError)
	}

	log.WithField("claims", claims).Debug("access token claims")
}

func (mx *Mux) expired(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) health(w http.ResponseWriter, req *http.Request) {
	payload := struct {
		Message    string `json:"message"`
		Status     string `json:"status"`
		StatusCode int    `json:"code"`
	}{
		fmt.Sprintf("application health status OK - %s", time.Now()),
		http.StatusText(http.StatusOK),
		http.StatusOK,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Add chi middleware to set content type on resp headers
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (mx *Mux) login(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) logout(w http.ResponseWriter, req *http.Request) {

}

func (mx *Mux) token(w http.ResponseWriter, req *http.Request) {

}
