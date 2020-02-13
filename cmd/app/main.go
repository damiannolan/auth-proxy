package main

import (
	"net/http"

	"github.com/damiannolan/auth-proxy/auth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{"port": viper.GetString("services.auth-proxy.port"), "server": viper.GetString("services.auth-proxy.host")}).Info("starting application server")

	authMux := auth.NewMux("/oauth")
	http.ListenAndServe(":"+viper.GetString("proxy.port"), authMux.Handler())
}
