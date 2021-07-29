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

func span(arr []float64) float64 {
	min := arr[0]
	max := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

func stepsFromOne(v float64) int {
	step := 0
	if v > 1 {
		for v > 1 {
			step++
			v /= 10
		}
	} else if v < 1 {
		for v < 1 {
			step--
			v *= 10
		}
	}
	return step
}

func decimals(arr ...[]float64) int {
	dp := 0

	for _, vals := range arr {
		s := stepsFromOne(span(vals))
		d := 0
		if s < 1 {
			d = -s + 1
		} else if s == 1 {
			d = 1
		}
		if dp < d {
			dp = d
		}
	}
	return dp
}
