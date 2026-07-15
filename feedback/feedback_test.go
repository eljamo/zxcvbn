package feedback_test

import (
	"testing"

	"github.com/eljamo/zxcvbn/feedback"
	"github.com/eljamo/zxcvbn/match"
	"github.com/stretchr/testify/assert"
)

const extraSuggestion = "Add another word or two. Uncommon words are better"

func TestGetFeedbackEmptySequence(t *testing.T) {
	fb := feedback.GetFeedback(0, nil)

	assert.Equal(t, "", fb.Warning)
	assert.Equal(t, []string{
		"Use a few words, avoid common phrases",
		"No need for symbols, digits, or uppercase letters",
	}, fb.Suggestions)
}

func TestGetFeedbackScoreAboveThreshold(t *testing.T) {
	seq := []*match.Match{
		{Pattern: "dictionary", DictionaryName: "passwords", Rank: 5, Token: "password"},
	}

	for _, score := range []int{3, 4} {
		fb := feedback.GetFeedback(score, seq)
		assert.Equal(t, feedback.Feedback{}, fb)
	}
}

func TestGetFeedbackLongestMatchWins(t *testing.T) {
	seq := []*match.Match{
		{Pattern: "dictionary", DictionaryName: "passwords", Rank: 5, Token: "pass"},
		{Pattern: "dictionary", DictionaryName: "english_wikipedia", Rank: 100, Token: "computer"},
	}

	fb := feedback.GetFeedback(0, seq)

	// The longer token ("computer") belongs to the second match, but it is not
	// a sole match, so english_wikipedia yields no warning.
	assert.Equal(t, "", fb.Warning)
	assert.Equal(t, extraSuggestion, fb.Suggestions[len(fb.Suggestions)-1])
}

func TestGetFeedbackDictionaryWarnings(t *testing.T) {
	tests := []struct {
		name    string
		match   *match.Match
		warning string
	}{
		{
			name:    "top-10 password",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "passwords", Rank: 5, Token: "password"},
			warning: "This is a top-10 common password",
		},
		{
			name:    "top-100 password",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "passwords", Rank: 50, Token: "password"},
			warning: "This is a top-100 common password",
		},
		{
			name:    "very common password",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "passwords", Rank: 500, Token: "password"},
			warning: "This is a very common password",
		},
		{
			name:    "l33t password low guesses",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "passwords", Rank: 500, L33t: true, Guesses: 100, Token: "p@5$w0rd"},
			warning: "This is similar to a commonly used password",
		},
		{
			name:    "sole wikipedia word",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "english_wikipedia", Rank: 100, Token: "computer"},
			warning: "A word by itself is easy to guess",
		},
		{
			name:    "sole female name",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "female_names", Rank: 10, Token: "jessica"},
			warning: "Names and surnames by themselves are easy to guess",
		},
		{
			name:    "sole male name",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "male_names", Rank: 10, Token: "james"},
			warning: "Names and surnames by themselves are easy to guess",
		},
		{
			name:    "sole surname",
			match:   &match.Match{Pattern: "dictionary", DictionaryName: "surnames", Rank: 10, Token: "smith"},
			warning: "Names and surnames by themselves are easy to guess",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := feedback.GetFeedback(0, []*match.Match{tt.match})
			assert.Equal(t, tt.warning, fb.Warning)
		})
	}
}

func TestGetFeedbackDictionarySuggestions(t *testing.T) {
	tests := []struct {
		name       string
		match      *match.Match
		suggestion string
	}{
		{
			name:       "l33t substitution",
			match:      &match.Match{Pattern: "dictionary", DictionaryName: "passwords", Rank: 500, L33t: true, Guesses: 100, Token: "p@5$w0rd"},
			suggestion: "Predictable substitutions like '@' instead of 'a' don't help very much",
		},
		{
			name:       "capitalized token",
			match:      &match.Match{Pattern: "dictionary", DictionaryName: "english_wikipedia", Rank: 100, Token: "Computer"},
			suggestion: "Capitalization doesn't help very much",
		},
		{
			name:       "all-caps token",
			match:      &match.Match{Pattern: "dictionary", DictionaryName: "english_wikipedia", Rank: 100, Token: "COMPUTER"},
			suggestion: "All-uppercase is almost as easy to guess as all-lowercase",
		},
		{
			name:       "reversed token",
			match:      &match.Match{Pattern: "dictionary", DictionaryName: "english_wikipedia", Rank: 100, Reversed: true, Token: "retupmoc"},
			suggestion: "Reversed words aren't much harder to guess",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := feedback.GetFeedback(0, []*match.Match{tt.match})
			assert.Contains(t, fb.Suggestions, tt.suggestion)
		})
	}
}

// TestGetFeedbackAllUpperDigitDivergence pins a deliberate divergence from
// dropbox/zxcvbn: this fork's AllUpper regex excludes digits (`^[^a-z\d]+$` vs
// upstream `^[^a-z]+$`), so digit-containing all-caps tokens do not get the
// all-uppercase suggestion. Do not "fix" the regex to match upstream without
// also inverting this test.
func TestGetFeedbackAllUpperDigitDivergence(t *testing.T) {
	m := &match.Match{Pattern: "dictionary", DictionaryName: "english_wikipedia", Rank: 100, Token: "PASSWORD123"}

	fb := feedback.GetFeedback(0, []*match.Match{m})

	assert.NotContains(t, fb.Suggestions, "All-uppercase is almost as easy to guess as all-lowercase")
}

func TestGetFeedbackNameInMultiMatchSequence(t *testing.T) {
	seq := []*match.Match{
		{Pattern: "dictionary", DictionaryName: "passwords", Rank: 5, Token: "ab"},
		{Pattern: "dictionary", DictionaryName: "female_names", Rank: 10, Token: "jessica"},
	}

	fb := feedback.GetFeedback(0, seq)

	assert.Equal(t, "Common names and surnames are easy to guess", fb.Warning)
}

// TestGetFeedbackRepeatWarningUsesBaseToken guards the regression where the
// warning was chosen by RepeatCount == 1 — a value the repeat matcher never
// produces (it is always >= 2) — instead of by base-token length.
func TestGetFeedbackRepeatWarningUsesBaseToken(t *testing.T) {
	// the match the repeat matcher actually produces for "aaaaaa"
	seq := []*match.Match{
		{Pattern: "repeat", BaseToken: "a", RepeatCount: 6, Token: "aaaaaa", I: 0, J: 5},
	}

	fb := feedback.GetFeedback(0, seq)

	assert.Equal(t, `Repeats like "aaa" are easy to guess`, fb.Warning)
}

func TestGetFeedbackNonDictionaryPatterns(t *testing.T) {
	tests := []struct {
		name    string
		match   *match.Match
		warning string
	}{
		{
			name:    "spatial straight row",
			match:   &match.Match{Pattern: "spatial", Turns: 1, Token: "qwerty"},
			warning: "Straight rows of keys are easy to guess",
		},
		{
			name:    "spatial with turns",
			match:   &match.Match{Pattern: "spatial", Turns: 3, Token: "qwaszx"},
			warning: "Short keyboard patterns are easy to guess",
		},
		{
			name:    "single-char base repeat",
			match:   &match.Match{Pattern: "repeat", BaseToken: "a", RepeatCount: 3, Token: "aaa"},
			warning: `Repeats like "aaa" are easy to guess`,
		},
		{
			name:    "multi-char base repeat",
			match:   &match.Match{Pattern: "repeat", BaseToken: "abc", RepeatCount: 3, Token: "abcabcabc"},
			warning: `Repeats like "abcabcabc" are only slightly harder to guess than "abc"`,
		},
		{
			name:    "multibyte single-rune base repeat",
			match:   &match.Match{Pattern: "repeat", BaseToken: "ä", RepeatCount: 3, Token: "äää"},
			warning: `Repeats like "aaa" are easy to guess`,
		},
		{
			name:    "sequence",
			match:   &match.Match{Pattern: "sequence", Token: "abcdef"},
			warning: `Sequences like "abc" or "6543" are easy to guess.`,
		},
		{
			name:    "recent year",
			match:   &match.Match{Pattern: "regex", RegexName: "recent_year", Token: "2025"},
			warning: "Recent years are easy to guess.",
		},
		{
			name:    "other regex",
			match:   &match.Match{Pattern: "regex", RegexName: "digits", Token: "123456"},
			warning: "",
		},
		{
			name:    "date",
			match:   &match.Match{Pattern: "date", Token: "13.8.1991"},
			warning: "Dates are often easy to guess.",
		},
		{
			name:    "unknown pattern",
			match:   &match.Match{Pattern: "bruteforce", Token: "zx!!"},
			warning: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := feedback.GetFeedback(0, []*match.Match{tt.match})
			assert.Equal(t, tt.warning, fb.Warning)
			// GetFeedback always appends the extra suggestion.
			assert.Equal(t, extraSuggestion, fb.Suggestions[len(fb.Suggestions)-1])
		})
	}
}
