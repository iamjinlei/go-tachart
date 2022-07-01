package charts

import (
	"log"

	"github.com/sosnovski/go-tachart/datasets"
	"github.com/sosnovski/go-tachart/opts"
	"github.com/sosnovski/go-tachart/render"
	"github.com/sosnovski/go-tachart/types"
)

// Geo represents a geo chart.
type Geo struct {
	BaseConfiguration
}

// Type returns the chart type.
func (Geo) Type() string { return types.ChartGeo }

var geoFormatter = `function (params) {
		return params.name + ' : ' + params.value[2];
}`

// NewGeo creates a new geo chart.
func NewGeo() *Geo {
	c := &Geo{}
	c.initBaseConfiguration()
	c.Renderer = render.NewChartRender(c, c.Validate)
	c.hasGeo = true
	return c
}

// AddSeries adds new data sets.
// geoType options:
// * types.ChartScatter
// * types.ChartEffectScatter
// * types.ChartHeatMap
func (c *Geo) AddSeries(name, geoType string, data []opts.GeoData, options ...SeriesOpts) *Geo {
	series := SingleSeries{Name: name, Type: geoType, Data: data, CoordSystem: types.ChartGeo}
	series.configureSeriesOpts(options...)
	c.MultiSeries = append(c.MultiSeries, series)
	return c
}

func (c *Geo) extendValue(region string, v float32) []float32 {
	res := make([]float32, 0)
	tv := datasets.Coordinates[region]
	if tv == [2]float32{0, 0} {
		log.Printf("goecharts: No coordinate is specified for %s\n", region)
	} else {
		res = append(tv[:], v)
	}
	return res
}

// SetGlobalOptions sets options for the Geo instance.
func (c *Geo) SetGlobalOptions(options ...GlobalOpts) *Geo {
	c.BaseConfiguration.setBaseGlobalOptions(options...)
	return c
}

// Validate
func (c *Geo) Validate() {
	if c.Tooltip.Formatter == "" {
		c.Tooltip.Formatter = opts.FuncOpts(geoFormatter)
	}
	c.Assets.Validate(c.AssetsHost)
}
