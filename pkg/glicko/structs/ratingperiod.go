// Reference to variables described on the specification can be found by
// comments with `doc-ref`.

package structs

import (
	"math"
)

// RatingPeriod holds information about the period to which the Glicko2
// calculation will be based on.
type RatingPeriod struct {
	SystemConstant float64
	Matches        []*Match
	Competitors    []*Competitor
}

// BuildRatingPeriod build a default RatingPeriod.
func BuildRatingPeriod() *RatingPeriod {
	return &RatingPeriod{
		SystemConstant: 0.5, // doc-ref: τ
		Matches:        []*Match{},
	}
}

// AddBuiltMatch adds Matches "instances" to the RatingPeriod.
// While adding the Match to the RatingPeriod, already register the other
// necessary data (link the Match also to the Competitors and "extract"
// the Competitors to register on the RatingPeriod).
func (rt *RatingPeriod) AddBuiltMatch(match *Match) {
	rt.Matches = append(rt.Matches, match)
	rt.addNewCompetitors(match.Home, match.Away)
	match.Home.AddMatch(match)
	match.Away.AddMatch(match)
}

// AddNewMatch adds creates Matches and add them to the RatingPeriod.
// While adding the Match to the RatingPeriod, already register the other
// necessary data (link the Match also to the Competitors and "extract"
// the Competitors to register on the RatingPeriod).
func (rt *RatingPeriod) AddNewMatch(home *Competitor, away *Competitor, winner int) {
	match := &Match{
		Home:   home,
		Away:   away,
		Winner: winner,
	}

	rt.Matches = append(rt.Matches, match)
	rt.addNewCompetitors(home, away)
	home.AddMatch(match)
	away.AddMatch(match)
}

// Calculate is responsible to glue all the magic together. At the end of it
// all the Competitors have the `.PostRating` data, which contains the new
// Glicko2 information after the results of the RatingPeriod.
func (rt *RatingPeriod) Calculate() {
	for _, competitor := range rt.Competitors {
		newVolatility := newVolatility(competitor, rt.SystemConstant)
		v := v(competitor)

		newPreRatingDerivation := math.Sqrt(math.Pow(competitor.PreRating.G2RatingDerivation, 2) + math.Pow(newVolatility, 2)) // doc-ref: φ*

		agg := 0.0
		for _, match := range competitor.Matches {
			opponent := match.OpponentOf(competitor)

			g := g(opponent.PreRating.G2RatingDerivation)
			E := e(competitor.PreRating.G2Rating, opponent.PreRating.G2Rating, opponent.PreRating.G2RatingDerivation)
			result := match.CompetitorResult(competitor)

			agg += g * (result - E)
		}

		// Each attribute is set in one individual line instead of constructing
		// the struct because each one depend on the result of the previous.
		competitor.PostRating = &Rating{}
		competitor.PostRating.G2RatingDerivation = 1 / (math.Sqrt((1 / math.Pow(newPreRatingDerivation, 2)) + (1 / v)))            // doc-ref: φ'
		competitor.PostRating.G2Rating = competitor.PreRating.G2Rating + math.Pow(competitor.PostRating.G2RatingDerivation, 2)*agg // doc-ref: µ'
		competitor.PostRating.Rating = 173.7178*competitor.PostRating.G2Rating + 1500                                              // doc-ref: r'
		competitor.PostRating.RatingDerivation = 173.7178 * competitor.PostRating.G2RatingDerivation                               // doc-ref: RD'
		competitor.PostRating.Volatility = newVolatility                                                                           // doc-ref: σ'
	}
}

// addCompetitor adds a Competitor to the list of Competitors inside the
// RatingPeriod.
func (rt *RatingPeriod) addCompetitor(competitor *Competitor) {
	rt.Competitors = append(rt.Competitors, competitor)
}

// containsCompetitor checks if the Competitor is already registered on the
// list if Competitors of the RatingPeriod.
func (rt *RatingPeriod) containsCompetitor(competitor *Competitor) bool {
	for _, c := range rt.Competitors {
		if c.ID == competitor.ID {
			return true
		}
	}
	return false
}

// addNewCompetitors adds a "set" Competitor to the list of Competitors
// inside the RatingPeriod.
func (rt *RatingPeriod) addNewCompetitors(competitors ...*Competitor) {
	for _, competitor := range competitors {
		if !rt.containsCompetitor(competitor) {
			rt.addCompetitor(competitor)
		}
	}
}

// Functions "translated" from the Glicko2 specification.
// The methods bellow this line don't contain documentation since they are
// pretty straightforward code-wise and any attempt of document the actual
// formulas would be a repetition of the specification. If any questions
// related to them occur, please check the glicko2.pdf located on this project.

func v(competitor *Competitor) float64 {
	agg := 0.0
	for _, match := range competitor.Matches {
		opponent := match.OpponentOf(competitor)

		g := g(opponent.PreRating.G2RatingDerivation)
		E := e(competitor.PreRating.G2Rating, opponent.PreRating.G2Rating, opponent.PreRating.G2RatingDerivation)

		agg += math.Pow(g, 2) * E * (1 - E)
	}

	return math.Pow(agg, -1)
}

func g(glicko2RatingDerivation float64) float64 {
	return 1 / (math.Sqrt(1 + (3*math.Pow(glicko2RatingDerivation, 2))/math.Pow(math.Pi, 2)))
}

func e(baseCompetitorGlicko2Rating, opponentGlicko2Rating, opponentGlicko2RatingDerivation float64) float64 {
	g := g(opponentGlicko2RatingDerivation)
	return 1 / (1 + math.Exp(-g*(baseCompetitorGlicko2Rating-opponentGlicko2Rating)))
}

func delta(competitor *Competitor) float64 {
	v := v(competitor)
	agg := 0.0
	for _, match := range competitor.Matches {
		opponent := match.OpponentOf(competitor)

		g := g(opponent.PreRating.G2RatingDerivation)
		E := e(competitor.PreRating.G2Rating, opponent.PreRating.G2Rating, opponent.PreRating.G2RatingDerivation)
		result := match.CompetitorResult(competitor)

		agg += g * (result - E)
	}

	return v * agg
}

func a(competitor *Competitor) float64 {
	return math.Log(math.Pow(competitor.PreRating.Volatility, 2))
}

func f(x float64, competitor *Competitor, constant float64) float64 {
	deltaPow := math.Pow(delta(competitor), 2)
	g2RatingDerivationPow := math.Pow(competitor.PreRating.G2RatingDerivation, 2)
	ePow := math.Pow(math.E, x)
	v := v(competitor)

	topLeft := ePow * (deltaPow - g2RatingDerivationPow - v - ePow)
	bottomLeft := 2 * math.Pow(g2RatingDerivationPow+v+ePow, 2)
	topRight := x - a(competitor)
	bottomRight := math.Pow(constant, 2)

	return (topLeft / bottomLeft) - (topRight / bottomRight)
}

// constant: doc-ref: τ
func newVolatility(competitor *Competitor, constant float64) float64 {
	var B float64
	A := a(competitor)
	v := v(competitor)
	delta := delta(competitor)
	e := 0.000001 // doc-ref: ε

	if math.Pow(delta, 2) > math.Pow(competitor.PreRating.G2RatingDerivation, 2)+v {
		B = math.Log(math.Pow(delta, 2) - math.Pow(competitor.PreRating.G2RatingDerivation, 2) - v)
	} else {
		k := 1.0
		for {
			x := A - (k * constant)
			if f(x, competitor, constant) < 0 {
				k += 1.0
			} else {
				B = A - (k * constant)
				break
			}
		}
	}

	fA := f(A, competitor, constant)
	fB := f(B, competitor, constant)

	for {
		if math.Abs(B-A) > e {
			C := A + ((A - B) * fA / (fB - fA))
			fC := f(C, competitor, constant)

			if fC*fB < 0 {
				A = B
				fA = fB
			} else {
				fA = fA / 2
			}

			B = C
			fB = fC
		} else {
			return math.Pow(math.E, (A / 2))
		}
	}
}
