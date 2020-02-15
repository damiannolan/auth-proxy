package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/damiannolan/auth-proxy/auth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {
	if err := readConfig(); err != nil {
		panic(err)
	}

	authMux := auth.NewMux("/oauth")

	host, port := viper.GetString("services.auth-proxy.host"), viper.GetInt("services.auth-proxy.port")
	log.WithFields(log.Fields{"host": host, "port": port}).Info("starting application server")
	http.ListenAndServe(fmt.Sprintf(":%d", port), authMux.Handler())
}
