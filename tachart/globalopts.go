package tachart

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/sosnovski/go-tachart/charts"
	"github.com/sosnovski/go-tachart/opts"
)

type globalOptsData struct {
	w           int64
	init        opts.Initialization
	titles      []opts.Title
	tooltip     opts.Tooltip
	axisPointer opts.AxisPointer
	grids       []opts.Grid
	xAxis       opts.XAxis // candlestick+overlay
	yAxis       opts.YAxis // candlestick+overlay
	dataZooms   []opts.DataZoom
}

func (c globalOptsData) genOpts(cfg Config, n int, eventDescMap map[string]string) []charts.GlobalOpts {
	tooltip := c.tooltip
	tooltip.Formatter = strings.Replace(tooltip.Formatter, "__EVENT_MAP__", toJson(eventDescMap), 1)

	return []charts.GlobalOpts{
		charts.WithTitleOpts(c.titles...),
		charts.WithInitializationOpts(c.init),
		charts.WithTooltipOpts(tooltip),
		charts.WithAxisPointerOpts(c.axisPointer),
		charts.WithGridOpts(c.grids...),
		charts.WithXAxisOpts(c.xAxis),
		charts.WithYAxisOpts(c.yAxis),
		charts.WithDataZoomOpts(c.dataZooms...),
	}
}

func toJson(o interface{}) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(o)
	return string(buf.Bytes())
}
