package matching

import (
	"testing"

	"github.com/eljamo/zxcvbn/match"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var testl33tTable = map[string][]string{
	"a": {"4", "@"},
	"c": {"(", "{", "[", "<"},
	"g": {"6", "9"},
	"o": {"0"},
	"s": {"5", "$"},
}

func Test_l33tMatch(t *testing.T) {
	lm := l33tMatch{
		dm: dictionaryMatch{
			rankedDictionaries: map[string]rankedDictionnary{
				"words": {
					"aac":       1,
					"password":  3,
					"paassword": 4,
					"asdf0":     5,
				},
				"words2": {
					"cgo": 1,
				},
			},
		},
		table: testl33tTable,
	}
	tests := []struct {
		name     string
		password string
		want     []*match.Match
	}{
		{
			name:     "doesn't match ''",
			password: "",
			want:     []*match.Match{},
		},
		{
			name:     "doesn't match pure dictionary words",
			password: "password",
			want:     []*match.Match{},
		},
		{
			name:     "matches against common l33t substitutions",
			password: "p4ssword",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "p4ssword",
					MatchedWord:    "password",
					Rank:           3,
					DictionaryName: "words",
					I:              0,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"4": "a"},
				},
			},
		},
		{
			name:     "matches against common l33t substitutions",
			password: "p@ssw0rd",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "p@ssw0rd",
					MatchedWord:    "password",
					Rank:           3,
					DictionaryName: "words",
					I:              0,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"@": "a", "0": "o"},
				},
			},
		},
		{
			name:     "matches against common l33t substitutions",
			password: "aSdfO{G0asDfO",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "{G0",
					MatchedWord:    "cgo",
					Rank:           1,
					DictionaryName: "words2",
					I:              5,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"{": "c", "0": "o"},
				},
			},
		},
		{
			name:     "matches against overlapping l33t patterns",
			password: "@a(go{G0",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "@a(",
					MatchedWord:    "aac",
					Rank:           1,
					DictionaryName: "words",
					I:              0,
					J:              2,
					L33t:           true,
					Sub:            map[string]string{"@": "a", "(": "c"},
				},
				{
					Pattern:        "dictionary",
					Token:          "(go",
					MatchedWord:    "cgo",
					Rank:           1,
					DictionaryName: "words2",
					I:              2,
					J:              4,
					L33t:           true,
					Sub:            map[string]string{"(": "c"},
				},
				{
					Pattern:        "dictionary",
					Token:          "{G0",
					MatchedWord:    "cgo",
					Rank:           1,
					DictionaryName: "words2",
					I:              5,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"{": "c", "0": "o"},
				},
			},
		},
		{
			name:     "matches when multiple l33t substitutions are needed for the same letter",
			password: "p@5$word",
			want: []*match.Match{
				{
					Pattern:        "dictionary",
					Token:          "p@5$word",
					MatchedWord:    "password",
					Rank:           3,
					DictionaryName: "words",
					I:              0,
					J:              7,
					L33t:           true,
					Sub:            map[string]string{"@": "a", "5": "s", "$": "s"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lm.Matches(tt.password))
		})
	}

	// doesn't match single-character l33ted words
	assert.Len(t, lm.Matches("4 1 @"), 0)

	assert.Equal(t, []*match.Match{
		{
			Pattern:        "dictionary",
			Token:          "4sdf0",
			MatchedWord:    "asdf0",
			Rank:           5,
			DictionaryName: "words",
			I:              0,
			J:              4,
			L33t:           true,
			Sub:            map[string]string{"4": "a"},
		},
	}, lm.Matches("4sdf0"))
}

func TestLeetTrieDeterministicOutput(t *testing.T) {
	password := "coRrecth0rseba++ery9.23.2007staple$"

	lm := l33tMatch{
		dm:    defaultRankedDictionaries,
		table: l33tTable,
	}

	var lastMatches []*match.Match
	for i := range 100 {
		matches := lm.Matches(password)
		if i > 0 {
			if d := cmp.Diff(matches, lastMatches); d != "" {
				t.Fatalf("Got two different values %s %s \n%s", match.ToString(lastMatches), match.ToString(matches), d)
			}
		}
		lastMatches = matches
	}
}
