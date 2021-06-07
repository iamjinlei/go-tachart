package tachart

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

type globalOptsData struct {
	init        opts.Initialization
	tooltip     opts.Tooltip
	axisPointer opts.AxisPointer
	grids       []opts.Grid
	xAxis       opts.XAxis // candlestick+overlay
	yAxis       opts.YAxis // candlestick+overlay
	dataZoom    opts.DataZoom
}

func (c globalOptsData) genOpts(eventDescMap map[string]string) []charts.GlobalOpts {
	tooltip := c.tooltip
	tooltip.Formatter = strings.Replace(tooltip.Formatter, "__EVENT_MAP__", toJson(eventDescMap), 1)

	return []charts.GlobalOpts{
		charts.WithInitializationOpts(c.init),
		charts.WithTooltipOpts(tooltip),
		charts.WithAxisPointerOpts(c.axisPointer),
		charts.WithGridOpts(c.grids...),
		charts.WithXAxisOpts(c.xAxis),
		charts.WithYAxisOpts(c.yAxis),
		charts.WithDataZoomOpts(c.dataZoom),
	}
}

func toJson(o interface{}) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(o)
	return string(buf.Bytes())
}
