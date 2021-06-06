package tachart

type OverlayType string

const (
	SMA   OverlayType = "SMA"
	EMA   OverlayType = "EMA"
	FIXED OverlayType = "FIXED"
)

type OverlayChart struct {
	Type OverlayType
	// SMA, EMA
	N int
	// FIXED
	Values []float64
}
