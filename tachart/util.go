package tachart

import (
	"fmt"
	"strings"
)

func countDecimalPlaces(v float64) int {
	// This is dumb, but math.Floor() introduces floating err
	parts := strings.Split(fmt.Sprintf("%v", v), ".")
	if len(parts) == 1 {
		return 0
	}
	return len(parts[1])
}

func maxDecimalPlaces(cdls []Candle) int {
	max := 0
	for _, cdl := range cdls {
		n := countDecimalPlaces(cdl.L)
		if n > max {
			max = n
		}
	}
	return max
}
