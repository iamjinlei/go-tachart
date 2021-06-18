package tachart

import (
	"fmt"
	"strings"

	"github.com/iamjinlei/go-tart"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

const (
	chartLabelFontSize   = 11
	chartLabelFontHeight = 13
)

type Indicator interface {
	name() string
	yAxisLabel() string
	yAxisMin() string
	yAxisMax() string
	getTitleOpts(top, left int, color string) []opts.Title
	genChart(vals []float64, xAxis interface{}, gridIndex int, color string) charts.Overlaper
}

type ma struct {
	nm string
	n  int64
	fn func([]float64, int64) []float64
}

func NewSMA(n int) Indicator {
	return ma{
		nm: fmt.Sprintf("SMA(%v)", n),
		n:  int64(n),
		fn: tart.SmaArr,
	}
}

func NewEMA(n int) Indicator {
	return ma{
		nm: fmt.Sprintf("EMA(%v)", n),
		n:  int64(n),
		fn: tart.EmaArr,
	}
}

func (c ma) name() string {
	return c.nm
}

func (c ma) yAxisLabel() string {
	return ""
}

func (c ma) yAxisMin() string {
	return ""
}

func (c ma) yAxisMax() string {
	return ""
}

func (c ma) getTitleOpts(top, left int, color string) []opts.Title {
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    color,
				FontSize: chartLabelFontSize,
			},
			Title: c.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (c ma) genChart(vals []float64, xAxis interface{}, gridIndex int, color string) charts.Overlaper {
	ma := c.fn(vals, c.n)
	for i := 0; i < int(c.n); i++ {
		ma[i] = ma[c.n]
	}

	items := []opts.LineData{}
	for _, v := range ma {
		items = append(items, opts.LineData{Value: v})
	}

	if color == "" {
		color = lineColors[0]
	}

	return charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm, items,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color: color,
			}))
}

type macd struct {
	nm     string
	fast   int64
	slow   int64
	signal int64
}

func NewMACD(fast, slow, signal int) Indicator {
	return macd{
		nm:     fmt.Sprintf("MACD(%v,%v,%v)", fast, slow, signal),
		fast:   int64(fast),
		slow:   int64(slow),
		signal: int64(signal),
	}
}

func (c macd) name() string {
	return c.nm
}

func (c macd) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (c macd) yAxisMin() string {
	return strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (c macd) yAxisMax() string {
	return strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (c macd) getTitleOpts(top, left int, _ string) []opts.Title {
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    lineColors[0],
				FontSize: chartLabelFontSize,
			},
			Title: c.nm + "-Diff",
			Left:  px(left),
			Top:   px(top),
		},
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    lineColors[1],
				FontSize: chartLabelFontSize,
			},
			Title: c.nm + "-Sig",
			Left:  px(left),
			Top:   px(top + chartLabelFontHeight),
		},
	}
}

func (c macd) genChart(vals []float64, xAxis interface{}, gridIndex int, _ string) charts.Overlaper {
	macd, signal, hist := tart.MacdArr(vals, c.fast, c.slow, c.signal)

	lineItems := []opts.LineData{}
	for _, v := range macd {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}
	macdLine := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Diff", lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color: lineColors[0],
			}),
		)

	lineItems = []opts.LineData{}
	for _, v := range signal {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}
	signalLine := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Sig", lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color: lineColors[1],
			}),
		)

	barItems := []opts.BarData{}
	for _, v := range hist {
		style := &opts.ItemStyle{
			Color:   colorUpBar,
			Opacity: opacity,
		}
		if v < 0 {
			style = &opts.ItemStyle{
				Color:   colorDownBar,
				Opacity: opacity,
			}
		}
		barItems = append(barItems, opts.BarData{Value: v, ItemStyle: style})
	}
	histBar := charts.NewBar().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Hist", barItems, charts.WithBarChartOpts(opts.BarChart{
			XAxisIndex: gridIndex,
			YAxisIndex: gridIndex,
			ZLevel:     100,
		}))

	macdLine.Overlap(signalLine, histBar)

	return macdLine
}
