package tachart

import (
	"fmt"

	"github.com/iamjinlei/go-tart"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

type atr struct {
	nm string
	n  int64
}

func NewATR(n int) Indicator {
	return atr{
		nm: fmt.Sprintf("ATR(%v)", n),
		n:  int64(n),
	}
}

func (a atr) name() string {
	return a.nm
}

func (a atr) yAxisLabel() string {
	return ""
}

func (a atr) yAxisMin() string {
	return ""
}

func (a atr) yAxisMax() string {
	return ""
}

func (a atr) getTitleOpts(top, left int, color string) []opts.Title {
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    color,
				FontSize: chartLabelFontSize,
			},
			Title: a.nm,
			Left:  px(left),
			Top:   px(top),
		},
	}
}

func (a atr) genChart(_, highs, lows, closes, _ []float64, xAxis interface{}, gridIndex int, color string) charts.Overlaper {
	vols := tart.AtrArr(highs, lows, closes, a.n)
	for i := 0; i < int(a.n); i++ {
		vols[i] = vols[a.n]
	}

	items := []opts.LineData{}
	for _, v := range vols {
		items = append(items, opts.LineData{Value: v})
	}

	if color == "" {
		color = lineColors[0]
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
				Color:   color,
				Opacity: opacityMed,
			}))
}
