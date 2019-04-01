package glicko

// RankableCompetitor is the struct that holds the data related to Glicko for a
// competitor.
type RankableCompetitor struct {
	ID         int
	PreRating  *Rating
	PostRating *Rating
	Matches    []*RankableMatch
}

// BuildRankableCompetitor is a builder to build a Competitor based on its id and a
// rating.
func BuildRankableCompetitor(id int, rating *Rating) *RankableCompetitor {
	return &RankableCompetitor{
		ID:        id,
		PreRating: rating,
	}
}

// AddMatch is used to add Mathes for a Competitor on a specific context
func (c *RankableCompetitor) AddMatch(match *RankableMatch) {
	c.Matches = append(c.Matches, match)
}
