package matching

import (
	"strconv"
	"testing"

	"github.com/eljamo/zxcvbn/match"
	"github.com/eljamo/zxcvbn/scoring"
	"github.com/stretchr/testify/assert"
)

func TestRegexpMatching(t *testing.T) {
	rm := regexpMatch{regexes: defaultRegexpMatch}
	recentYear := strconv.Itoa(scoring.ReferenceYear)
	assert.Equal(t, []*match.Match{
		{
			Pattern:   "regex",
			Token:     recentYear,
			I:         0,
			J:         3,
			RegexName: "recent_year",
		},
	},
		rm.Matches(recentYear),
	)

	futureYear := strconv.Itoa(scoring.ReferenceYear + 1)
	assert.Equal(t, []*match.Match{
		{
			Pattern:   "regex",
			Token:     futureYear,
			I:         0,
			J:         3,
			RegexName: "recent_year",
		},
	},
		rm.Matches(futureYear),
	)

	assert.Empty(t, rm.Matches(strconv.Itoa(scoring.ReferenceYear+scoring.MinYearSpace+1)))
	assert.Empty(t, rm.Matches(strconv.Itoa(scoring.ReferenceYear-recentYearPastWindow-1)))
}
