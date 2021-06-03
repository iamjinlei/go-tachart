package tachart

type OverlayType string

const (
	SMA   OverlayType = "sma"
	EMA   OverlayType = "ema"
	FIXED OverlayType = "fixed"
)

type OverlayChart struct {
	Type OverlayType
	// SMA, EMA
	N int
	// FIXED
	Values []float64
}
