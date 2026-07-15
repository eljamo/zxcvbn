package zxcvbn

import (
	"math"
	"strings"
	"testing"
	"unicode/utf8"
)

// PasswordStrength's match search is O(n^2) in the password length by design
// (upstream dropbox/zxcvbn is the same), so unbounded fuzz inputs eventually
// exceed the fuzzer's per-exec time budget without indicating a bug. This
// bound is far above any real password length while keeping the worst
// observed case (repeated characters) well under a second.
const maxFuzzPasswordLen = 1024

// FuzzPasswordStrength enforces the library-wide match invariants: every
// emitted match must satisfy Token == password[I:J+1] with I and J+1 on rune
// boundaries (byte-index convention — do not "fix" the assertions to rune
// indexing). The seed corpus doubles as deterministic regression coverage for
// the multibyte index fixes (l33t and sequence matchers) and runs on every
// ordinary `go test ./...` invocation.
func FuzzPasswordStrength(f *testing.F) {
	seeds := []string{
		"",
		"password",
		"p@ssword",
		"äp@ssword",      // regression: l33t rune indices
		"üp4sswordé",     // regression: l33t rune indices
		"стуфхц",         // regression: byte-based sequences (U+0441..U+0446)
		"стувгд",         // two adjacent sequence runs (сту + вгд)
		"abcд",           // regression: rune-splitting tokens
		"kwyjibokwyjibo", // repeat matcher with base analysis
		"ääää",           // multibyte repeat
		"13/8/1991",      // date
		"qwertyuio",      // spatial
		"correcthorsebatterystaple",
		"\xff\xfe", // invalid UTF-8: must not panic, returns zero Result
		strings.Repeat("U", maxFuzzPasswordLen), // worst-case length: repeated chars maximize match counts
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, password string) {
		if len(password) > maxFuzzPasswordLen {
			return
		}

		result := PasswordStrength(password, nil)

		if !utf8.ValidString(password) {
			// documented early-return path: zero Result
			return
		}

		for _, m := range result.Sequence {
			if m.I < 0 || m.J < m.I || m.J >= len(password) {
				t.Fatalf("match %q (%s): indices i=%d j=%d out of range for password %q",
					m.Token, m.Pattern, m.I, m.J, password)
			}
			if got := password[m.I : m.J+1]; got != m.Token {
				t.Fatalf("match (%s): token %q != password[%d:%d] %q in %q",
					m.Pattern, m.Token, m.I, m.J+1, got, password)
			}
			if !utf8.ValidString(m.Token) {
				t.Fatalf("match (%s): token %q is not valid UTF-8 (password %q)",
					m.Pattern, m.Token, password)
			}
		}

		if len(password) > 0 {
			if math.IsNaN(result.Guesses) || math.IsInf(result.Guesses, 0) || result.Guesses < 1 {
				t.Fatalf("guesses=%v for password %q", result.Guesses, password)
			}
			if result.Score < 0 || result.Score > 4 {
				t.Fatalf("score=%d for password %q", result.Score, password)
			}
		}
	})
}
