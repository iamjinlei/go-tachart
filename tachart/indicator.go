package tachart

import (
	"github.com/sosnovski/go-tachart/charts"
	"github.com/sosnovski/go-tachart/opts"
)

const (
	chartLabelFontSize   = 11
	chartLabelFontHeight = 13
)

type Indicator interface {
	// Name indicator name
	Name() string
	// YAxisLabel label formatter
	YAxisLabel() string
	// YAxisMin label formatter
	YAxisMin() string
	// YAxisMax label formatter
	YAxisMax() string
	// GetNumColors # of colors needed
	GetNumColors() int
	// GetTitleOpts indicator chart legend config
	GetTitleOpts(top, left int, colorIndex int) []opts.Title
	// GenChart indicator chart config
	GenChart(opens, highs, lows, closes, vols []float64, xAxis interface{}, gridIndex int) charts.Overlaper
}
