package main

import (
	"fmt"
	"notifier/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// TODO: SSL disabled bad?
func initDb() *sqlx.DB {
	log.Debug().Msg(fmt.Sprintf("postgresql://postgres:%s@%s/canigowatchthis", config.Config.PsqlPassword, config.Config.PsqlUrl))
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgresql://postgres:%s@%s/canigowatchthis?sslmode=disable", config.Config.PsqlPassword, config.Config.PsqlUrl))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init database")
	}

	var i int
	r, err := db.Query("SELECT COUNT(*) FROM games;")
	if err != nil {
		log.Error().Err(err).Msg("failed to count")
		return db
	}
	if r.Next() {
		r.Scan(&i)
		log.Debug().Int("i", i).Msg("")
	}

	return db
}
