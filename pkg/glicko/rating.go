package glicko

// Rating is the struct that holds the Glicko2 data.
type Rating struct {
	Rating             float64 // doc-ref: r
	RatingDerivation   float64 // doc-ref: RD
	Volatility         float64 // doc-ref: σ
	G2Rating           float64 // doc-ref: µ
	G2RatingDerivation float64 // doc-ref: φ
}

// BuildRating build a Rating using positional args.
func BuildRating(args ...float64) *Rating {
	rating := &Rating{
		Rating:           args[0],
		RatingDerivation: args[1],
		Volatility:       args[2],
	}

	// While building the ratings, already calculate the Glicko2 equivalents.
	rating.G2Rating = (rating.Rating - 1500) / 173.7178
	rating.G2RatingDerivation = rating.RatingDerivation / 173.7178

	return rating
}

// BuildDefaultRating creates a Rating with default values.
func BuildDefaultRating() *Rating {
	return BuildRating(1500, 350, 0.06)
}
