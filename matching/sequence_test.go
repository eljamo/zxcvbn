package matching

import (
	"testing"
	"unicode/utf8"

	"github.com/eljamo/zxcvbn/match"
	"github.com/stretchr/testify/assert"
)

func Test_sequenceMatch_Matches(t *testing.T) {
	s := sequenceMatch{}

	// doesn't match 1- and 2-character spatial patterns
	assert.Empty(t, s.Matches(""))
	assert.Empty(t, s.Matches("a"))
	assert.Empty(t, s.Matches("1"))

	// matches overlapping patterns
	assert.Equal(t, []*match.Match{
		{
			Pattern:       "sequence",
			Token:         "abc",
			I:             0,
			J:             2,
			Ascending:     true,
			SequenceName:  "lower",
			SequenceSpace: 26,
		},
		{
			Pattern:       "sequence",
			Token:         "cba",
			I:             2,
			J:             4,
			Ascending:     false,
			SequenceName:  "lower",
			SequenceSpace: 26,
		},
		{
			Pattern:       "sequence",
			Token:         "abc",
			I:             4,
			J:             6,
			Ascending:     true,
			SequenceName:  "lower",
			SequenceSpace: 26,
		},
	}, s.Matches("abcbabc"))

	// matches embedded sequence patterns
	word := "jihg"
	for _, pv := range genpws(word, []string{"!", "22"}, []string{"!", "22"}) {
		assert.Equal(t, []*match.Match{
			{
				Pattern:       "sequence",
				Token:         word,
				I:             pv.i,
				J:             pv.j,
				Ascending:     false,
				SequenceName:  "lower",
				SequenceSpace: 26,
			},
		}, s.Matches(pv.password))
	}

	// doesn't match a single multibyte rune (the old byte-based length guard
	// would have entered the loop)
	assert.Empty(t, s.Matches("ä"))

	// matches sequences of consecutive code points (U+0441..U+0446), with
	// byte-offset I/J
	assert.Equal(t, []*match.Match{
		{
			Pattern:       "sequence",
			Token:         "стуфхц",
			I:             0,
			J:             11,
			Ascending:     true,
			SequenceName:  "unicode",
			SequenceSpace: 26,
		},
	}, s.Matches("стуфхц"))

	// multibyte runes don't split or shift ASCII sequence indices
	mixed := "abcд"
	mixedMatches := s.Matches(mixed)
	assert.Equal(t, []*match.Match{
		{
			Pattern:       "sequence",
			Token:         "abc",
			I:             0,
			J:             2,
			Ascending:     true,
			SequenceName:  "lower",
			SequenceSpace: 26,
		},
	}, mixedMatches)
	for _, m := range mixedMatches {
		assert.Equal(t, mixed[m.I:m.J+1], m.Token)
		assert.True(t, utf8.ValidString(m.Token))
	}

	// matches pattern with the right sequence type
	tests := []struct {
		pattern   string
		name      string
		ascending bool
		space     int
	}{
		{"ABC", "upper", true, 26},
		{"CBA", "upper", false, 26},
		{"PQR", "upper", true, 26},
		{"RQP", "upper", false, 26},
		{"XYZ", "upper", true, 26},
		{"ZYX", "upper", false, 26},
		{"abcd", "lower", true, 26},
		{"dcba", "lower", false, 26},
		{"jihg", "lower", false, 26},
		{"wxyz", "lower", true, 26},
		{"zxvt", "lower", false, 26},
		{"0369", "digits", true, 10},
		{"97531", "digits", false, 10},
	}
	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			matches := s.Matches(tt.pattern)
			assert.Equal(t, []*match.Match{{
				Pattern:       "sequence",
				Token:         tt.pattern,
				I:             0,
				J:             len(tt.pattern) - 1,
				Ascending:     tt.ascending,
				SequenceName:  tt.name,
				SequenceSpace: tt.space,
			}}, matches)
		})
	}
}
