package tachart

import (
	"fmt"
	"strings"

	"github.com/iamjinlei/go-tart"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

type rsi struct {
	nm         string
	n          int64
	oversold   float64
	overbought float64
}

func NewRSI(n int, oversold, overbought float64) Indicator {
	return rsi{
		nm:         fmt.Sprintf("RSI(%v)", n),
		n:          int64(n),
		oversold:   oversold,
		overbought: overbought,
	}
}

func (r rsi) name() string {
	return r.nm
}

func (r rsi) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (r rsi) yAxisMin() string {
	return `function(value) { return 0 }`
}

func (r rsi) yAxisMax() string {
	return `function(value) { return 100 }`
}

func (r rsi) getTitleOpts(top, left int, _ string) []opts.Title {
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    lineColors[0],
				FontSize: chartLabelFontSize,
			},
			Title: r.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (r rsi) genChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int, _ string) charts.Overlaper {
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
				Color:   lineColors[0],
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
