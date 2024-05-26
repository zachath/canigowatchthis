package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	PsqlPassword string `mapstructure:"postgres_password"`
	PsqlUrl      string `mapstructure:"postgres_url"`
}

const (
	psqlPassword = ""
	psqlUrl      = "postgres.canigowatchthis"
)

var Config AppConfig

func LoadConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("postgres_password", psqlPassword)
	viper.SetDefault("postgres_url", psqlUrl)

	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load config")
	}
}
