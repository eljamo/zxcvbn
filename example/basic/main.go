package main

import (
	"fmt"

	"github.com/eljamo/zxcvbn"
)

func main() {
	password := "p@5$w0rd"

	result := zxcvbn.PasswordStrength(password, nil)

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Score: %d/4\n", result.Score)
	fmt.Printf("Guesses: %.0f\n", result.Guesses)
	fmt.Printf("Guesses (log10): %.2f\n", result.GuessesLog10)
	fmt.Printf("Online throttled crack time: %s\n", result.CrackTimesDisplay["online_throttling_100_per_hour"])
	fmt.Printf("Online unthrottled crack time: %s\n", result.CrackTimesDisplay["online_no_throttling_10_per_second"])
	fmt.Printf("Offline fast hash crack time: %s\n", result.CrackTimesDisplay["offline_fast_hashing_1e10_per_second"])
	fmt.Printf("Offline slow hash crack time: %s\n", result.CrackTimesDisplay["offline_slow_hashing_1e4_per_second"])
	fmt.Printf("Calculation time: %.3fs\n", result.CalcTime)
}
