package main

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func readConfig() error {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

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

	viper.WatchConfig()

	return nil
}
