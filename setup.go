package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfiguration(filename string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(filename)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	return v
}

func LoadLogger(cfg *viper.Viper) logrus.FieldLogger {
	l := logrus.New()

	lvl, err := logrus.ParseLevel(cfg.GetString("app.log_level"))
	if err != nil {
		panic(err)
	}

	l.SetLevel(lvl)
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetOutput(os.Stderr)

	return l
}
