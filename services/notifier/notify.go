package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Game struct {
	Id        int       `db:"id"`
	GameId    string    `db:"game_id"`
	League    string    `db:"league"`
	StartTime time.Time `db:"start_time"`
}

func notify() error {
	log.Debug().Msg("I am not implemented")
	return nil
}
