package tachart

import (
	"strings"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

type line struct {
	nm   string
	vals []float64
}

func NewLine(name string, vals []float64) Indicator {
	return line{
		nm:   name,
		vals: vals,
	}
}

func (b line) name() string {
	return b.nm
}

func (b line) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (b line) yAxisMin() string {
	return ""
}

func (b line) yAxisMax() string {
	return ""
}

func (b line) getTitleOpts(top, left int, _ string) []opts.Title {
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    lineColors[0],
				FontSize: chartLabelFontSize,
			},
			Title: b.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (b line) genChart(_, _, _, _, _ []float64, xAxis interface{}, gridIndex int, _ string) charts.Overlaper {
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
				Color:   lineColors[0],
				Opacity: opacityMed,
			}),
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
