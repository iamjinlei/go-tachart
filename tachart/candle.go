package tachart

type Candle struct {
	Label string  // x-axis label for this candle. Usually a timestamp, e.g. "2018/1/24 08:00"
	O     float64 // open
	H     float64 // high
	L     float64 // low
	C     float64 // close
	V     float64 // volume
}
