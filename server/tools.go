package server

import (
	"time"
)

func randomizeColor() string {
	randSeed := time.Now().UnixNano()
	randSeed = randSeed % int64(230) // int64 for more randomness due to UnixNano() being int64
	return "\033[38;5;" + itoa(int(randSeed)) + "m"
}

func itoa(n int) string {
	result := ""
	isNegative := ""

	if n < 0 {
		n *= 1
		isNegative = "-"
	}

	if n == 0 {
		return result
	}

	for n > 0 {
		result = string(rune(n%10+'0')) + result
		n /= 10
	}

	return isNegative + result
}
