package main

import (
	"api/config"
	"api/nhlapi"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

type InputContext int

const (
	location InputContext = iota
	spans
	team
)

func InputValidation(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userInput UserInput
		err := json.NewDecoder(r.Body).Decode(&userInput)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to decode input")
			JSONResponse(w, "failed to parse input data", http.StatusBadRequest)
			return
		}

		if !nhlapi.IsValidTeam(strings.ToLower(userInput.Team)) {
			log.Error().Str("team", userInput.Team).Msg("invalid team provided")
			JSONResponse(w, fmt.Sprintf("invalid team '%s' provided", userInput.Team), http.StatusBadRequest)
			return
		}

		input := Input{
			Timezone: userInput.Timezone,
			Spans:    map[Weekday]Span{},
			Team:     strings.ToLower(userInput.Team),
		}

		for _, s := range maps.Keys(userInput.Spans) {
			input.Spans[weekDayFromString(s)] = userInput.Spans[s]
		}

		loc, err := time.LoadLocation(input.Timezone)
		if err != nil {
			log.Error().Stack().Err(err).Str("timezone", input.Timezone).Msg("failed parse timezeon")
			JSONResponse(w, fmt.Sprintf("'%s' is an invalid timezone", input.Timezone), http.StatusBadRequest)
			return
		}

		if len(input.Spans) == 0 {
			log.Error().Stack().Err(err).Msg("no spans provided")
			JSONResponse(w, "empty list of spans provided", http.StatusBadRequest)
			return
		}

		for i, s := range input.Spans {
			if i < Any || i > Sunday {
				log.Error().Stack().Err(err).Int("weekday", int(i)).Msg("invalid weekday provided")
				JSONResponse(w, fmt.Sprintf("invalid weekday '%d'", i), http.StatusBadRequest)
				return
			}

			if s.Start < 0 || s.Start > 24 {
				log.Error().Stack().Err(err).Int("hour", s.Start).Msg("invalid start hour provided")
				JSONResponse(w, fmt.Sprintf("invalid start hour '%d' provied in span [%d]", s.Start, i), http.StatusBadRequest)
				return
			}

			if s.End < 0 || s.End > 24 {
				log.Error().Stack().Err(err).Int("hour", s.End).Msg("invalid end hour provided")
				JSONResponse(w, fmt.Sprintf("invalid end hour '%d' provied in span [%d]", s.End, i), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, location, loc)
		ctx = context.WithValue(ctx, spans, input.Spans)
		ctx = context.WithValue(ctx, team, input.Team)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ProcessRequest(w http.ResponseWriter, r *http.Request) {
	location := r.Context().Value(location).(*time.Location)
	spans := r.Context().Value(spans).(map[Weekday]Span)
	team := r.Context().Value(team).(string)

	requests.WithLabelValues(team).Inc()

	schedule, err := nhlapi.GetTeamSchedule(team)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get team schedule")
		JSONResponse(w, fmt.Sprintf("failed to get schedule of team '%s'", team), http.StatusBadRequest)
		return
	}

	response := ResponseMessage{
		Games: []nhlapi.Game{},
	}

	for _, game := range schedule.Games {
		gameDateUTC, err := time.Parse(time.RFC3339, game.StartTimeUTC)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("failed to parse %s", game.StartTimeUTC)
			continue
		}

		if config.Config.IgnoreOldGames && gameDateUTC.Before(time.Now().UTC()) {
			continue
		}

		gameDateUserTimezone := gameDateUTC.In(location)
		weekday := convertWeekday(gameDateUserTimezone.Weekday())

		// If the specific weekday has not been provided nor the 'Any' weekday, skip. Otherwise use the 'Any' span.
		if !slices.Contains(maps.Keys(spans), weekday) {
			if !slices.Contains(maps.Keys(spans), Any) {
				continue
			}

			weekday = Any
		}

		span := spans[weekday]

		lower := time.Date(
			gameDateUserTimezone.Year(),
			gameDateUserTimezone.Month(),
			gameDateUserTimezone.Day(),
			span.Start,
			0,
			0,
			0,
			location,
		)

		higher := time.Date(
			gameDateUserTimezone.Year(),
			gameDateUserTimezone.Month(),
			gameDateUserTimezone.Day(),
			span.End,
			0,
			0,
			0,
			location,
		)

		if inSpan(gameDateUserTimezone, lower, higher) {
			response.Games = append(response.Games, game)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	if len(response.Games) == 0 {
		status = http.StatusNotFound
		response.Message = "No games found for provided spans"
	}

	w.WriteHeader(status)
	response.StatusCode = status
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Stack().Err(err).Interface("message", response).Msg("failed to encode error message")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"statusCode": 500, "message": "Internal Server Error"}`))
	}
}

func inSpan(gameTime, lower, higher time.Time) bool {
	return (gameTime.After(lower) && gameTime.Before(higher)) || gameTime.Equal(lower) || gameTime.Equal(higher)
}

type ResponseMessage struct {
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Games      []nhlapi.Game `json:"games,omitempty"`
}

func JSONResponse(w http.ResponseWriter, msg string, code int) {
	errorsCount.WithLabelValues(http.StatusText(code)).Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := ResponseMessage{
		StatusCode: code,
		Message:    msg,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Stack().Err(err).Interface("message", response).Msg("failed to encode error message")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"statusCode": 500, "message": "Internal Server Error"}`))
	}
}
