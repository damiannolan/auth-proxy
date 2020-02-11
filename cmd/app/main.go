package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello world...")

	http.HandleFunc("/oauth/authorize", nil)
	http.HandleFunc("/oauth/callback", nil)
	http.HandleFunc("/oauth/expired", nil)
	http.HandleFunc("/oauth/health", nil)
	http.HandleFunc("/oauth/login", nil)
	http.HandleFunc("/oauth/logout", nil)
	http.HandleFunc("/oauth/authorize", nil)
	http.HandleFunc("/oauth/token", nil)

	http.ListenAndServe(":8080", nil)
}
