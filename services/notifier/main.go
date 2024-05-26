package main

import (
	"notifier/config"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zachath/canigowatchthis/libs/pingcaplog"
)

var (
	db *sqlx.DB
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.ErrorStackMarshaler = pingcaplog.PingcapMarshalStack

	config.LoadConfig()
	db = initDb()
}

func main() {
	log.Info().Msg("starting notification job")
	err := notify()
	if err != nil {
		log.Fatal().Err(err).Msg("notification job failed")
	}
}
