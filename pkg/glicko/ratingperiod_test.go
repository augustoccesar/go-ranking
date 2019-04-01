package glicko

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockRatingPeriod(t *testing.T) *RatingPeriod {
	competitor1 := BuildRankableCompetitor(1, BuildRating(1500, 200, 0.06))
	competitor2 := BuildRankableCompetitor(2, BuildRating(1400, 30, 0.06))
	competitor3 := BuildRankableCompetitor(3, BuildRating(1550, 100, 0.06))
	competitor4 := BuildRankableCompetitor(4, BuildRating(1700, 300, 0.06))

	ratingPeriod := BuildRatingPeriod(1)

	ratingPeriod.AddNewMatch(competitor1, competitor2, 1)
	ratingPeriod.AddNewMatch(competitor1, competitor3, 3)
	ratingPeriod.AddNewMatch(competitor1, competitor4, 4)

	return ratingPeriod
}

func TestAddNewMatch(t *testing.T) {
	competitor1 := BuildRankableCompetitor(1, BuildRating(0.0, 0.0, 0.0))
	competitor2 := BuildRankableCompetitor(2, BuildRating(0.0, 0.0, 0.0))

	ratingPeriod := BuildRatingPeriod(1)

	ratingPeriod.AddNewMatch(competitor1, competitor2, 1)

	assert.Equal(t, 2, len(ratingPeriod.Competitors))
	assert.Equal(t, 1, len(ratingPeriod.Matches))
}

func TestG(t *testing.T) {
	ratingPeriod := mockRatingPeriod(t)
	competitor1 := ratingPeriod.Competitors[0]
	competitor2 := ratingPeriod.Competitors[1]
	competitor3 := ratingPeriod.Competitors[2]
	competitor4 := ratingPeriod.Competitors[3]

	result1 := g(competitor1.PreRating.G2RatingDerivation)
	result2 := g(competitor2.PreRating.G2RatingDerivation)
	result3 := g(competitor3.PreRating.G2RatingDerivation)
	result4 := g(competitor4.PreRating.G2RatingDerivation)

	assert.LessOrEqual(t, math.Abs(0.8442-result1), 0.0001)
	assert.LessOrEqual(t, math.Abs(0.9955-result2), 0.0001)
	assert.LessOrEqual(t, math.Abs(0.9532-result3), 0.0001)
	assert.LessOrEqual(t, math.Abs(0.7242-result4), 0.0001)
}

func TestE(t *testing.T) {
	ratingPeriod := mockRatingPeriod(t)
	competitor1 := ratingPeriod.Competitors[0]
	competitor2 := ratingPeriod.Competitors[1]
	competitor3 := ratingPeriod.Competitors[2]
	competitor4 := ratingPeriod.Competitors[3]

	result1 := e(competitor1.PreRating.G2Rating, competitor2.PreRating.G2Rating, competitor2.PreRating.G2RatingDerivation)
	result2 := e(competitor1.PreRating.G2Rating, competitor3.PreRating.G2Rating, competitor3.PreRating.G2RatingDerivation)
	result3 := e(competitor1.PreRating.G2Rating, competitor4.PreRating.G2Rating, competitor4.PreRating.G2RatingDerivation)

	assert.LessOrEqual(t, math.Abs(0.639-result1), 0.001)
	assert.LessOrEqual(t, math.Abs(0.432-result2), 0.001)
	assert.LessOrEqual(t, math.Abs(0.303-result3), 0.001)
}

func TestV(t *testing.T) {
	ratingPeriod := mockRatingPeriod(t)
	competitor1 := ratingPeriod.Competitors[0]

	result := v(competitor1)

	assert.LessOrEqual(t, math.Abs(1.7785-result), 0.0005)
}

func TestDelta(t *testing.T) {
	ratingPeriod := mockRatingPeriod(t)
	competitor1 := ratingPeriod.Competitors[0]

	result := delta(competitor1)

	assert.LessOrEqual(t, math.Abs(-0.4834-result), 0.001)
}

func TestNewVolatile(t *testing.T) {
	ratingPeriod := mockRatingPeriod(t)
	competitor1 := ratingPeriod.Competitors[0]

	result := newVolatility(competitor1, ratingPeriod.SystemConstant)

	assert.LessOrEqual(t, math.Abs(0.05999-result), 0.00001)
}

func TestCalculate(t *testing.T) {
	ratingPeriod := mockRatingPeriod(t)
	ratingPeriod.Calculate()

	competitor1 := ratingPeriod.Competitors[0]

	assert.LessOrEqual(t, math.Abs(1464.06-competitor1.PostRating.Rating), 0.1)
	assert.LessOrEqual(t, math.Abs(151.52-competitor1.PostRating.RatingDerivation), 0.1)
	assert.LessOrEqual(t, math.Abs(0.05999-competitor1.PostRating.Volatility), 0.00001)
}
