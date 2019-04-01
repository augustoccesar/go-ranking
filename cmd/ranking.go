package main

import (
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/augustoccesar/go-ranking/internal/spider/thescore"
	"github.com/augustoccesar/go-ranking/pkg/glicko"
	"github.com/urfave/cli"
)

// AvailableSources contains the list of sources that the script supports
var AvailableSources = []string{"thescore"}

func checkSource(source string) bool {
	for _, availabeSource := range AvailableSources {
		if availabeSource == source {
			return true
		}
	}
	return false
}

type InputParams struct {
	Source         string
	StartDate      string
	EndDate        string
	PeriodDuration int64
}

func main() {
	inputParams := InputParams{}

	app := cli.NewApp()
	app.Name = "Ranking CLI."
	app.Usage = ""

	defaltEndTime := time.Now().UTC()
	defaultStartTime := defaltEndTime.AddDate(0, -1, 0)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "source",
			Value:       "thescore",
			Usage:       "Source from where the system will fetch data.",
			Destination: &inputParams.Source,
		},
		cli.StringFlag{
			Name:        "start_date",
			Value:       defaultStartTime.Format(time.RFC3339),
			Usage:       "Date that the system will use as base to look for data.",
			Destination: &inputParams.StartDate,
		},
		cli.StringFlag{
			Name:        "end_date",
			Value:       defaltEndTime.Format(time.RFC3339),
			Usage:       "Limit date of the data.",
			Destination: &inputParams.EndDate,
		},
		cli.Int64Flag{
			Name:        "period_duration",
			Value:       7,
			Usage:       "Length in days of the Rating Period.",
			Destination: &inputParams.PeriodDuration,
		},
	}

	app.Action = func(c *cli.Context) error {
		parsedStartDate, _ := time.Parse(time.RFC3339, inputParams.StartDate)
		parsedEndDate, _ := time.Parse(time.RFC3339, inputParams.EndDate)
		switch inputParams.Source {
		case "thescore":
			daysGap := parsedEndDate.Sub(parsedStartDate).Hours() / 24
			periods := math.Ceil(daysGap / float64(inputParams.PeriodDuration))
			periodData, _ := thescore.FetchPeriodData(parsedStartDate, parsedEndDate)

			teamsCache := map[int]*thescore.Team{}
			teamsRating := map[int]*glicko.Rating{}

			ratingPeriods := []*glicko.RatingPeriod{}
			for i := 0; i < int(periods); i++ {
				ratingPeriod := glicko.BuildRatingPeriod(i + 1)

				sort.SliceStable(periodData.Matches, func(i, j int) bool {
					return periodData.Matches[i].StartTime.Before(periodData.Matches[j].StartTime)
				})

				periodCompetitors := map[int]*glicko.RankableCompetitor{}

				for _, match := range periodData.Matches {
					var home *glicko.RankableCompetitor
					var away *glicko.RankableCompetitor
					var homeRating *glicko.Rating
					var awayRating *glicko.Rating

					if val, ok := teamsRating[match.Home.ID]; ok {
						homeRating = val
					} else {
						homeRating = glicko.BuildDefaultRating()
					}

					if val, ok := teamsRating[match.Away.ID]; ok {
						awayRating = val
					} else {
						awayRating = glicko.BuildDefaultRating()
					}

					if val, ok := periodCompetitors[match.Home.ID]; ok {
						home = val
					} else {
						home = &glicko.RankableCompetitor{
							ID:        match.Home.ID,
							PreRating: homeRating,
						}
						periodCompetitors[match.Home.ID] = home
					}

					if val, ok := periodCompetitors[match.Away.ID]; ok {
						away = val
					} else {
						away = &glicko.RankableCompetitor{
							ID:        match.Away.ID,
							PreRating: awayRating,
						}
						periodCompetitors[match.Home.ID] = away
					}

					// This can be improved
					teamsCache[match.Home.ID] = match.Home
					teamsCache[match.Away.ID] = match.Away

					var winnerID int
					if winner := match.Winner; winner != nil {
						winnerID = winner.ID
					} else {
						winnerID = -1
					}
					ratingPeriod.AddNewMatch(home, away, winnerID)
				}

				ratingPeriod.Calculate()

				for _, competitor := range ratingPeriod.Competitors {
					teamsRating[competitor.ID] = competitor.PostRating
				}

				ratingPeriods = append(ratingPeriods, ratingPeriod)

				log.Printf("Ranking by the end of period: %d\n\n", ratingPeriod.ID)
				sort.SliceStable(ratingPeriod.Competitors, func(i, j int) bool {
					return ratingPeriod.Competitors[i].PostRating.Rating > ratingPeriod.Competitors[j].PostRating.Rating
				})
				for i, competitor := range ratingPeriod.Competitors {
					variation := competitor.PostRating.Rating - competitor.PreRating.Rating
					variationSymbol := ""
					if variation > 0 {
						variationSymbol = "+"
					}

					log.Printf("\t#%d - %s - %f (%s%f)\n", i+1, teamsCache[competitor.ID].Name, competitor.PostRating.Rating, variationSymbol, variation)
				}
				log.Printf("------------------------------------------------\n")
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}