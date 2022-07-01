package tachart

import (
	"fmt"
	"strings"

	"github.com/sosnovski/go-tachart/charts"
	"github.com/sosnovski/go-tachart/opts"
)

type boundedLine struct {
	nm          string
	vals        []float64
	min         float64
	max         float64
	lowerMarker float64
	upperMarker float64
	ci          int
}

func NewBoundedLine(name string, vals []float64, min, max, lowerMarker, upperMarker float64) Indicator {
	return &boundedLine{
		nm:          name,
		vals:        vals,
		min:         min,
		max:         max,
		lowerMarker: lowerMarker,
		upperMarker: upperMarker,
	}
}

func (b boundedLine) Name() string {
	return b.nm
}

func (b boundedLine) YAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (b boundedLine) YAxisMin() string {
	return fmt.Sprintf("function(value) { return %v }", b.min)
}

func (b boundedLine) YAxisMax() string {
	return fmt.Sprintf("function(value) { return %v }", b.max)
}

func (b boundedLine) GetNumColors() int {
	return 1
}

func (b *boundedLine) GetTitleOpts(top, left int, colorIndex int) []opts.Title {
	b.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[b.ci],
				FontSize: chartLabelFontSize,
			},
			Title: b.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (b boundedLine) GenChart(_, _, _, _, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	lineItems := []opts.LineData{}
	for _, v := range b.vals {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}

	return charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm, lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[b.ci],
				Opacity: opacityMed,
			}),
			charts.WithMarkLineNameYAxisItemOpts(
				opts.MarkLineNameYAxisItem{
					Name:  "lower_bound ",
					YAxis: b.lowerMarker,
				},
				opts.MarkLineNameYAxisItem{
					Name:  "upper_bound",
					YAxis: b.upperMarker,
				},
			),
			charts.WithMarkLineStyleOpts(
				opts.MarkLineStyle{
					Symbol: []string{"none", "none"},
					LineStyle: &opts.LineStyle{
						Color:   colorDownBar,
						Opacity: opacityMed,
					},
				},
			),
		)
}
