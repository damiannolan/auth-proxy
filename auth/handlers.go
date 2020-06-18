package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/damiannolan/auth-proxy/openid"
	"github.com/damiannolan/auth-proxy/realm"
	log "github.com/sirupsen/logrus"
)

func (svc *ProxyService) authorize(w http.ResponseWriter, req *http.Request) {
	realmID, ok := realm.FromContext(req.Context())
	if !ok {
		redirectURL := "getRedirectURL()"
		http.Redirect(w, req, redirectURL, http.StatusTemporaryRedirect)
	}

	if p, ok := svc.providers[realmID]; ok { // "default"
		http.Redirect(w, req, p.Config.AuthCodeURL(req.URL.Query().Get("state")), http.StatusTemporaryRedirect)
	} else {
		discoveryURL := openid.BuildDiscoveryURL(realmID)
		authProvider := openid.NewAuthenticationProvider(discoveryURL)
		svc.providers[realmID] = authProvider

		http.Redirect(w, req, p.Config.AuthCodeURL(req.URL.Query().Get("state")), http.StatusTemporaryRedirect)
	}
}

func (svc *ProxyService) callback(w http.ResponseWriter, req *http.Request) {
	realmID, ok := realm.FromContext(req.Context())
	if !ok {
		http.Redirect(w, req, "", http.StatusTemporaryRedirect)
	}

	authProvider := svc.providers[realmID]

	code := req.URL.Query().Get("code")
	token, err := authProvider.Exchange(req.Context(), code)
	if err != nil {
		http.Redirect(w, req, "", http.StatusForbidden)
	}

	log.WithField("token", token).Debug("oauth2 token")

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Redirect(w, req, "", http.StatusForbidden)
	}

	idToken, err := authProvider.Verify(req.Context(), rawIDToken)
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

func (svc *ProxyService) expired(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement a function to check if token has expired
}

func (svc *ProxyService) health(w http.ResponseWriter, req *http.Request) {
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

func (svc *ProxyService) login(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement direct login handler
}

func (svc *ProxyService) logout(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement logout and revoke token
}

func (svc *ProxyService) token(w http.ResponseWriter, req *http.Request) {
	// TODO
}
