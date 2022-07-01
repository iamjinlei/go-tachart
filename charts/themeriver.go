package charts

import (
	"github.com/sosnovski/go-tachart/opts"
	"github.com/sosnovski/go-tachart/render"
	"github.com/sosnovski/go-tachart/types"
)

// ThemeRiver represents a theme river chart.
type ThemeRiver struct {
	BaseConfiguration
}

// Type returns the chart type.
func (ThemeRiver) Type() string { return types.ChartThemeRiver }

// NewThemeRiver creates a new theme river chart.
func NewThemeRiver() *ThemeRiver {
	c := &ThemeRiver{}
	c.initBaseConfiguration()
	c.Renderer = render.NewChartRender(c, c.Validate)
	c.hasSingleAxis = true
	return c
}

// AddSeries adds new data sets.
func (c *ThemeRiver) AddSeries(name string, data []opts.ThemeRiverData, options ...SeriesOpts) *ThemeRiver {
	cd := make([][3]interface{}, len(data))
	for i := 0; i < len(data); i++ {
		cd[i] = data[i].ToList()
	}
	series := SingleSeries{Name: name, Type: types.ChartThemeRiver, Data: cd}
	series.configureSeriesOpts(options...)
	c.MultiSeries = append(c.MultiSeries, series)
	return c
}

// SetGlobalOptions sets options for the ThemeRiver instance.
func (c *ThemeRiver) SetGlobalOptions(options ...GlobalOpts) *ThemeRiver {
	c.BaseConfiguration.setBaseGlobalOptions(options...)
	return c
}

// Validate
func (c *ThemeRiver) Validate() {
	c.Assets.Validate(c.AssetsHost)
}
