package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	LogPretty bool   `envconfig:"LOG_PRETTY" default:"false"`
	Port      int    `envconfig:"PORT" default:"8000"`
	Debug     bool   `envconfig:"DEBUG" default:"false"`

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

	initLogger(c)
	log.Info().Msg("Starting")
	log.Info().Int("Port", c.Port).Send()
	log.Info().Str("Log level", c.LogLevel).Send()
	log.Info().Bool("Log pretty", c.LogPretty).Send()
	log.Info().Str("Error Files path", c.ErrorFilesPath).Send()
	log.Info().Str("Default format", c.DefaultFormat).Send()

	return c
}
