package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	LogPretty bool   `envconfig:"LOG_PRETTY" default:"false"`
	Port      int    `envconfig:"PORT" default:"8000"`
}

func MustParseConfig() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		_ = envconfig.Usage("", &c)
		log.Fatal().Msg(err.Error())
	}

	return c
}
