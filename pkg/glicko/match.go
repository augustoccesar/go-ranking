package glicko

// RankableMatch is the struct that hold the information related to a Match.
type RankableMatch struct {
	Home   *RankableCompetitor
	Away   *RankableCompetitor
	Winner int
}

// BuildRankableMatch build a Match based on the Competitors and the ID of the Winner.
func BuildRankableMatch(home *RankableCompetitor, away *RankableCompetitor, winner int) *RankableMatch {
	return &RankableMatch{
		Home:   home,
		Away:   away,
		Winner: winner,
	}
}

// OpponentOf is used to find who is the opponent of a specific Competitor
// inside a Match.
func (m *RankableMatch) OpponentOf(competitor *RankableCompetitor) *RankableCompetitor {
	if m.Home.ID == competitor.ID {
		return m.Away
	}
	return m.Home
}

// CompetitorResult is used to get the result value that is expected by
// Glicko2 formulas based on the result of the Match.
func (m *RankableMatch) CompetitorResult(competitor *RankableCompetitor) float64 {
	if m.Winner == -1 {
		return 0.5 // Tie
	} else if m.Winner == competitor.ID {
		return 1
	} else {
		return 0
	}
}

// WinnerCompetitor is used to get the Competitor that won the Match.
func (m *RankableMatch) WinnerCompetitor() *RankableCompetitor {
	if m.Home.ID == m.Winner {
		return m.Home
	} else if m.Away.ID == m.Winner {
		return m.Away
	} else {
		return nil
	}
}
