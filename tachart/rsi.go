package tachart

import (
	"fmt"
	"strings"

	"github.com/sosnovski/go-tart"

	"github.com/sosnovski/go-tachart/charts"
	"github.com/sosnovski/go-tachart/opts"
)

type rsi struct {
	nm         string
	n          int64
	oversold   float64
	overbought float64
	ci         int
}

func NewRSI(n int, oversold, overbought float64) Indicator {
	return &rsi{
		nm:         fmt.Sprintf("RSI(%v)", n),
		n:          int64(n),
		oversold:   oversold,
		overbought: overbought,
	}
}

func (r rsi) Name() string {
	return r.nm
}

func (r rsi) YAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (r rsi) YAxisMin() string {
	return `function(value) { return 0 }`
}

func (r rsi) YAxisMax() string {
	return `function(value) { return 100 }`
}

func (r rsi) GetNumColors() int {
	return 1
}

func (r *rsi) GetTitleOpts(top, left int, colorIndex int) []opts.Title {
	r.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[r.ci],
				FontSize: chartLabelFontSize,
			},
			Title: r.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (r rsi) GenChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	vals := tart.RsiArr(closes, r.n)

	lineItems := []opts.LineData{}
	for _, v := range vals {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}

	return charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(r.nm, lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[r.ci],
				Opacity: opacityMed,
			}),
			charts.WithMarkLineNameYAxisItemOpts(
				opts.MarkLineNameYAxisItem{
					Name:  "oversold",
					YAxis: r.oversold,
				},
				opts.MarkLineNameYAxisItem{
					Name:  "overbought",
					YAxis: r.overbought,
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
