package thescore

import "time"

// PeriodData is the struct that will hold all the data related to the fetched
// data for a period.
type PeriodData struct {
	StartTime time.Time
	EndTime   time.Time
	Matches   []*Match
	Teams     []*Team

	matchesCache map[int]*Match
	teamsCache   map[int]*Team
}

// BuildPeriodData is used to build a PeriodData and call the necessary methods
// to ensure all necessary data during creation.
func BuildPeriodData(
	startTime time.Time, endTime time.Time,
	matches []*Match, teams []*Team,
) *PeriodData {
	periodData := &PeriodData{
		StartTime: startTime,
		EndTime:   endTime,
		Matches:   matches,
		Teams:     teams,
	}

	periodData.populateCache()
	periodData.assignTeamsToMatches()

	return periodData
}

// populateCache is used to load a small cache of the Teams and Matches on
// easy queryable maps.
func (pd *PeriodData) populateCache() {
	pd.matchesCache = map[int]*Match{}
	pd.teamsCache = map[int]*Team{}

	for _, match := range pd.Matches {
		pd.matchesCache[match.ID] = match
	}

	for _, team := range pd.Teams {
		pd.teamsCache[team.ID] = team
	}
}

// assignTeamsToMatches uses the ids found on the root of the match to query
// the cache and assing the Teams to the Matches
func (pd *PeriodData) assignTeamsToMatches() {
	for _, match := range pd.Matches {
		homeID, awayID, winnerID := match.extractTeamsIds()
		match.Home = pd.teamsCache[homeID]
		match.Away = pd.teamsCache[awayID]
		match.Winner = pd.teamsCache[winnerID]
	}
}
