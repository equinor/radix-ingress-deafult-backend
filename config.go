package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	LogPretty bool   `envconfig:"LOG_PRETTY" default:"false"`
	Port      int    `envconfig:"PORT" default:"8000"`

	ErrorFilesPath string `envconfig:"ERROR_FILES_PATH" default:"./www"`
	DefaultFormat  string `envconfig:"DEFAULT_RESPONSE_FORMAT" default:"text/html"`
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
