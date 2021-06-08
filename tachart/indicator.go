package tachart

import (
	"fmt"
	"strings"

	"github.com/markcheno/go-talib"

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
	n  int
	fn func([]float64, int) []float64
}

func NewSMA(n int) Indicator {
	return ma{
		nm: fmt.Sprintf("SMA(%v)", n),
		n:  n,
		fn: talib.Sma,
	}
}

func NewEMA(n int) Indicator {
	return ma{
		nm: fmt.Sprintf("EMA(%v)", n),
		n:  n,
		fn: talib.Ema,
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
	for i := 0; i < c.n; i++ {
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
	fast   int
	slow   int
	signal int
}

func NewMACD(fast, slow, signal int) Indicator {
	return macd{
		nm:     fmt.Sprintf("MACD(%v,%v,%v)", fast, slow, signal),
		fast:   fast,
		slow:   slow,
		signal: signal,
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
	macd, signal, hist := talib.Macd(vals, c.fast, c.slow, c.signal)

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
			Color: colorUpBar,
		}
		if v < 0 {
			style = &opts.ItemStyle{
				Color: colorDownBar,
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
