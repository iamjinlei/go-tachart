package tachart

import (
	"html/template"
	"path/filepath"
	"runtime"
)

// page is conceptually divided into 3x3 grids:
// ----------------------------------------------
//                      top
// ----------------------------------------------
//           |                       |
//    left   |          chart        |   right
//           |                       |
// ----------------------------------------------
//                      bottom
// ----------------------------------------------
type pageLayout struct {
	chartWidth    int
	chartHeight   int
	topContent    template.HTML
	topHeight     int
	bottomContent template.HTML
	bottomHeight  int
	leftContent   template.HTML
	leftWidth     int
	rightContent  template.HTML
	rightWidth    int
}

type Config struct {
	precision          int // decimal places of floating nubmers shown on chart
	overlays           []Indicator
	indicators         []Indicator
	assetsHost         string
	theme              Theme
	layout             pageLayout
	draggable          bool
	eventDescWrapWidth int // wrap width of event desc on tooltip, 0 means no-wrap
	jsFuncs            []string
	disableVol         bool
	showLegend         bool
}

func NewConfig() *Config {
	return &Config{
		precision:  2,
		overlays:   []Indicator{},
		indicators: []Indicator{},
		assetsHost: "https://go-echarts.github.io/go-echarts-assets/assets/",
		theme:      ThemeWhite,
		layout: pageLayout{
			chartWidth:  900,
			chartHeight: 500,
		},
		draggable:          false,
		eventDescWrapWidth: 160,
		disableVol:         false,
	}
}

func (c *Config) SetTheme(t Theme) *Config {
	c.theme = t
	return c
}

func (c *Config) SetChartWidth(w int) *Config {
	c.layout.chartWidth = w
	return c
}

func (c *Config) SetChartHeight(h int) *Config {
	c.layout.chartHeight = h
	return c
}

func (c *Config) SetTopRowContent(content string, h int) *Config {
	c.layout.topContent = template.HTML(content)
	c.layout.topHeight = h
	return c
}

func (c *Config) SetBottomRowContent(content string, h int) *Config {
	c.layout.bottomContent = template.HTML(content)
	c.layout.bottomHeight = h
	return c
}

func (c *Config) SetLeftColContent(content string, w int) *Config {
	c.layout.leftContent = template.HTML(content)
	c.layout.leftWidth = w
	return c
}

func (c *Config) SetRightColContent(content string, w int) *Config {
	c.layout.rightContent = template.HTML(content)
	c.layout.rightWidth = w
	return c
}

func (c *Config) SetDraggable(draggable bool) *Config {
	c.draggable = draggable
	return c
}

func (c *Config) SetDisableVols(disable bool) *Config {
	c.disableVol = disable
	return c
}

func (c *Config) SetEnableLegend(enabled bool) *Config {
	c.showLegend = enabled
	return c
}

func (c *Config) SetEventDescWrapWidth(w int) *Config {
	c.eventDescWrapWidth = w
	return c
}

func (c *Config) AddJSFunc(js string) *Config {
	c.jsFuncs = append(c.jsFuncs, js)
	return c
}

func (c *Config) SetPrecision(p int) *Config {
	c.precision = p
	return c
}

func (c *Config) AddOverlay(vals ...Indicator) *Config {
	c.overlays = append(c.overlays, vals...)
	return c
}

func (c *Config) AddIndicator(vals ...Indicator) *Config {
	c.indicators = append(c.indicators, vals...)
	return c
}

func (c *Config) UseRepoAssets() *Config {
	// serving assets from "this" repo in local file system
	// with accessing network
	_, path, _, _ := runtime.Caller(0)
	path = filepath.Dir(path)
	c.assetsHost = filepath.Join("file:/"+filepath.Dir(path), "assets/")
	return c
}

func (c *Config) SetAssetsHost(assetsHost string) *Config {
	// serving assets from specified host
	c.assetsHost = assetsHost
	return c
}
