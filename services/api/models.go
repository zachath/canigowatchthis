package main

import (
	"strings"
	"time"
)

type Weekday int

const (
	Any Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

var days = []string{"any", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

func weekDayFromString(s string) Weekday {
	lower := strings.ToLower(s)
	for i, d := range days {
		if d == lower {
			return Weekday(i)
		}
	}

	return Any
}

// I refuse to consider Sunday as the beginning of the week.
// Sets Sunday to equal 7
func convertWeekday(goWeekday time.Weekday) Weekday {
	if goWeekday == time.Sunday {
		return Sunday
	}

	return Weekday(goWeekday)
}

type UserInput struct {
	Timezone string          `json:"timezone"`
	Spans    map[string]Span `json:"spans"`
	Team     string          `json:"team"`
}

type Input struct {
	Timezone string           `json:"timezone"`
	Spans    map[Weekday]Span `json:"spans"`
	Team     string           `json:"team"`
}

type Span struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
