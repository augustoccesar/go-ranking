package thescore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// FetchPeriodData is used to get the PeriodData by a start and end time.
func FetchPeriodData(startTime, endTime time.Time) (*PeriodData, error) {
	// Parse the dates to the format expected by the API
	formatedStartTime := startTime.Format(time.RFC3339)
	formatedEndTime := endTime.Format(time.RFC3339)

	// Build the URL
	baseURL := "https://esports-api.thescore.com/csgo/matches?%s"
	filter := fmt.Sprintf("start_date_from=%s&start_date_to=%s",
		formatedStartTime,
		formatedEndTime,
	)
	url := fmt.Sprintf(baseURL, filter)

	// Handle possible error while getting the response
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Parse the JSON response
	var rootData map[string]*json.RawMessage
	var teams []*Team
	var matches []*Match

	json.Unmarshal(body, &rootData)
	json.Unmarshal(*rootData["teams"], &teams)
	json.Unmarshal(*rootData["matches"], &matches)

	return BuildPeriodData(startTime, endTime, matches, teams), nil
}
