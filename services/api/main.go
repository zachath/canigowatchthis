package main

import (
	"api/config"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zachath/canigowatchthis/libs/pingcaplog"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.ErrorStackMarshaler = pingcaplog.PingcapMarshalStack
}

func init() {
	config.LoadConfig()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/gui", nil) //TODO: rm
	r.HandleFunc("/api", InputValidation(ProcessRequest)).Methods(http.MethodPost)
	r.Handle("/metrics", promhttp.Handler()) // TODO: This should not be exposed externaly.

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", config.Config.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Info().Str("port", config.Config.Port).Msg("server started")
	log.Fatal().Err(server.ListenAndServe()).Msg("server exited")
}
