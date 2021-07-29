package tachart

import (
	"fmt"
	"strings"

	"github.com/iamjinlei/go-tart"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

type atr struct {
	nm string
	n  int64
	ci int
	dp int
}

func NewATR(n int) Indicator {
	return &atr{
		nm: fmt.Sprintf("ATR(%v)", n),
		n:  int64(n),
	}
}

func (a atr) name() string {
	return a.nm
}

func (a atr) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", a.dp), -1)
}

func (a atr) yAxisMin() string {
	return strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", a.dp), -1)
}

func (a atr) yAxisMax() string {
	return strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", a.dp), -1)
}

func (a atr) getNumColors() int {
	return 1
}

func (a *atr) getTitleOpts(top, left int, colorIndex int) []opts.Title {
	a.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[a.ci],
				FontSize: chartLabelFontSize,
			},
			Title: a.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (a atr) genChart(_, highs, lows, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	vals := tart.AtrArr(highs, lows, closes, a.n)
	for i := 0; i < int(a.n); i++ {
		vals[i] = vals[a.n]
	}
	a.dp = decimals(vals)

	items := []opts.LineData{}
	for _, v := range vals {
		items = append(items, opts.LineData{Value: v})
	}

	return charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(a.nm, items,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[a.ci],
				Opacity: opacityMed,
			}))
}
