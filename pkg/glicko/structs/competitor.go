package structs

// Competitor is the struct that holds the data related to Glicko for a
// competitor.
type Competitor struct {
	ID         int
	PreRating  *Rating
	PostRating *Rating
	Matches    []*Match
}

// BuildCompetitor is a builder to build a Competitor based on its id and a
// rating.
func BuildCompetitor(id int, rating *Rating) *Competitor {
	return &Competitor{
		ID:        id,
		PreRating: rating,
	}
}

// AddMatch is used to add Mathes for a Competitor on a specific context
func (c *Competitor) AddMatch(match *Match) {
	c.Matches = append(c.Matches, match)
}
