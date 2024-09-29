package nhlapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pingcap/errors"
)

type GameState string
type GameType int

const (
	baseUrl = "https://api-web.nhle.com/v1"

	GameStateFuture GameState = "FUT"
	GameStateFinal  GameState = "FINAL"

	Preseason GameType = 1
	Regular   GameType = 2
	Playoff   GameType = 3
)

var teams = [...]string{
	"det",
	"bos",
	"wpg",
	"pit",
	"tbl",
	"phi",
	"tor",
	"car",
	"uta",
	"cgy",
	"mtl",
	"wsh",
	"van",
	"col",
	"nsh",
	"ana",
	"vgk",
	"sea",
	"dal",
	"chi",
	"nyr",
	"cbj",
	"fla",
	"edm",
	"min",
	"stl",
	"ott",
	"nyi",
	"lak",
	"njd",
	"buf",
	"sjs",
}

type TeamSchedule struct {
	CurrentSeason int    `json:"currentSeason"`
	ClubUTCOffset string `json:"clubUTCOffset"`
	Games         []Game `json:"games"`
}

type Game struct {
	Type             GameType   `json:"gameType"`
	Date             string     `json:"gameDate"`
	Venue            NamedPlace `json:"venue"`
	StartTimeUTC     string     `json:"startTimeUTC"`
	EasternUTCOffset string     `json:"easternUTCOffset"`
	VenueUTCOffset   string     `json:"venueUTCOffset"`
	VenueTimezone    string     `json:"venueTimezone"`
	AwayTeam         Team       `json:"awayTeam"`
	HomeTeam         Team       `json:"homeTeam"`
}

type NamedPlace struct {
	Name string `json:"default"`
}

type Team struct {
	Id           int        `json:"id"`
	PlaceName    NamedPlace `json:"placeName"`
	Abbreviation string     `json:"abbrev"`
	Logo         string     `json:"logo"`
}

func GetTeamSchedule(team string) (TeamSchedule, error) {
	resp, err := http.Get(fmt.Sprintf("%s/club-schedule-season/%s/now", baseUrl, team))
	if err != nil {
		return TeamSchedule{}, errors.Annotatef(err, "failed to get schedule of team '%s'", team)
	}

	if resp.StatusCode != http.StatusOK {
		return TeamSchedule{}, errors.Annotatef(err, "received unexpected status code '%d'", resp.StatusCode)
	}

	var schedule TeamSchedule
	err = json.NewDecoder(resp.Body).Decode(&schedule)
	if err != nil {
		return TeamSchedule{}, errors.Annotate(err, "failed to decode response")
	}

	return schedule, nil
}

func IsValidTeam(team string) bool {
	for _, t := range teams {
		if t == team {
			return true
		}
	}
	return false
}
