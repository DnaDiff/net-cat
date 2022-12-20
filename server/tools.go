package server

import (
	"fmt"
	"strings"
	"time"
)

func toggleRoomName(name string) string {
	if strings.HasPrefix(name, "\033[36m") && strings.HasSuffix(name, "\033[0m") {
		return strings.TrimSuffix(strings.TrimPrefix(name, "\033[36m"), "\033[0m")
	} else if !strings.HasPrefix(name, "\033[36m") && !strings.HasSuffix(name, "\033[0m") {
		return "\033[36m" + name + "\033[0m"
	} else {
		fmt.Fprintf(mw, "Error: Could not toggle room name: %s", name)
		return name
	}
}

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
