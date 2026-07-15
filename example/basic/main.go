package main

import (
	"fmt"

	"github.com/eljamo/zxcvbn"
)

func main() {
	password := "p@5$w0rd"

	result := zxcvbn.PasswordStrength(password, nil)

	crackTimes := []struct {
		label string
		key   string
	}{
		{"Online throttled", "online_throttling_100_per_hour"},
		{"Online unthrottled", "online_no_throttling_10_per_second"},
		{"Offline slow hash", "offline_slow_hashing_1e4_per_second"},
		{"Offline fast hash", "offline_fast_hashing_1e10_per_second"},
	}

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Score: %d/4\n", result.Score)
	fmt.Printf("Guesses: %.0f (log10: %.2f)\n", result.Guesses, result.GuessesLog10)
	for _, ct := range crackTimes {
		fmt.Printf("%-20s %s\n", ct.label+":", result.CrackTimesDisplay[ct.key])
	}
	fmt.Printf("Feedback Warning: %s\n", result.Feedback.Warning)
	fmt.Printf("Feedback Suggestions: %v\n", result.Feedback.Suggestions)
	fmt.Printf("Calculation time: %.3fs\n", result.CalcTime)
}
