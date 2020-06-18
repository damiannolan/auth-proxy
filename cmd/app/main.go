package main

import (
	"os"
	"time"

	"github.com/damiannolan/auth-proxy/auth"
	log "github.com/sirupsen/logrus"
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
		log.WithError(err).Error("failed to read config")
	}

	srv, err := auth.NewProxyService()
	if err != nil {
		log.WithError(err).Error("failed to create proxy service")
	}

	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Error("failed to serve")
	}
}
