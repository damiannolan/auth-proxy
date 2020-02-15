package main

import (
	"net/http"
	"os"
	"time"

	"github.com/damiannolan/auth-proxy/auth"
	"github.com/fsnotify/fsnotify"
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
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.WithFields(log.Fields{"file": e.Name, "op": e.Op.String()}).Debug("Config file changed")

		lvl, err := log.ParseLevel(viper.GetString("logger.level"))
		if err != nil {
			log.WithError(err).WithField("level", log.DebugLevel).Info("setting default logger level")
			log.SetLevel(log.DebugLevel)
			return
		}

		log.WithField("level", viper.GetString("logger.level")).Info("setting logger level")
		log.SetLevel(lvl)
	})

	log.WithFields(log.Fields{"port": viper.GetString("services.auth-proxy.port"), "server": viper.GetString("services.auth-proxy.host")}).Info("starting application server")

	authMux := auth.NewMux("/oauth")
	http.ListenAndServe(":"+viper.GetString("proxy.port"), authMux.Handler())
}
