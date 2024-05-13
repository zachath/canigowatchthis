package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port           string `mapstructure:"port"`
	IgnoreOldGames bool   `mapstructure:"ignore_old_games"`
}

const (
	port           = "8080"
	ignoreOldGames = true
)

var Config AppConfig

func LoadConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("port", port)

	viper.SetDefault("ignore_old_games", ignoreOldGames)

	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load config")
	}
}
