package main

import (
	"net/http"

	"github.com/damiannolan/auth-proxy/auth"
)

func main() {
	authMux := auth.NewMux("/oauth")
	http.ListenAndServe(":8080", authMux.Handler())
}
