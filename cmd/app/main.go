package main

import (
	"net/http"

	"github.com/damiannolan/auth-proxy/auth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("application") // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("config")      // path to look for the config file in
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{"port": viper.GetString("proxy.port"), "server": viper.GetString("proxy.host")}).Info("starting application server")
	authMux := auth.NewMux("/oauth")
	http.ListenAndServe(":"+viper.GetString("proxy.port"), authMux.Handler())
}
