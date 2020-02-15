package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/damiannolan/auth-proxy/realm"
)

func (mx *Mux) authorize(w http.ResponseWriter, req *http.Request) {
	realmID, ok := realm.FromContext(req.Context())
	if !ok {
		// redirect
	}

	provider := mx.providers[realmID]
	http.Redirect(w, req, provider.cfg.AuthCodeURL(req.URL.Query().Get("state")), http.StatusTemporaryRedirect)
}

func (mx *Mux) callback(w http.ResponseWriter, req *http.Request) {
	realmID, _ := realm.FromContext(req.Context())
	provider := mx.providers[realmID]
	code := req.URL.Query().Get("code")

	token, err := provider.cfg.Exchange(req.Context(), code)
	if err != nil {
		// handle err
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		// handle missing token
	}

	// Parse and verify ID Token payload.
	idToken, err := provider.verifier.Verify(req.Context(), rawIDToken)
	if err != nil {
		// handle error
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		// handle error
	}
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
