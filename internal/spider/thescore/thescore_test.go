package thescore

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchPeriodData(t *testing.T) {
	// First Quarter-Finals of Katowice 2019
	startTime := time.Date(2019, 03, 01, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2019, 03, 02, 0, 0, 0, 0, time.UTC)

	periodData, _ := FetchPeriodData(startTime, endTime)

	assert.Equal(t, 2, len(periodData.Matches))
	assert.Equal(t, 4, len(periodData.Teams))

	assert.Equal(t, "MIBR", periodData.Teams[0].Name)

	fmt.Print(periodData)
}
