package thescore

import (
	"regexp"
	"strconv"
	"time"
)

// Match is the struct that represents TheScore API response for Match (with
// stripped down fields for only what I need).
type Match struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	HomeURL   string    `json:"team1_url"`
	AwayURL   string    `json:"team2_url"`
	HomeScore int       `json:"team1_score"`
	AwayScore int       `json:"team2_score"`
	TieMatch  bool      `json:"tie_match"`
	WinnerURL string    `json:"winning_team_url"`
	StartTime time.Time `json:"start_date"`

	Home   *Team
	Away   *Team
	Winner *Team
}

// extractTeamsIds uses regex to extract the Teams ids from the fields that
// consists of urls to the Teams
func (m *Match) extractTeamsIds() (homeID, awayID, winnerID int) {
	regex := *regexp.MustCompile(`\/csgo\/teams\/(\d+)`)

	homeResult := regex.FindStringSubmatch(m.HomeURL)
	awayResult := regex.FindStringSubmatch(m.AwayURL)
	winnerResult := regex.FindStringSubmatch(m.WinnerURL)

	homeID, _ = strconv.Atoi(homeResult[1])
	awayID, _ = strconv.Atoi(awayResult[1])
	if m.TieMatch {
		winnerID = -1
	} else {
		winnerID, _ = strconv.Atoi(winnerResult[1])
	}

	return homeID, awayID, winnerID
}
