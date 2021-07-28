package tachart

import (
	"fmt"

	"github.com/iamjinlei/go-tart"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

type bbands struct {
	nm      string
	n       int64
	nStdDev float64
	isSma   bool
	ci      int
}

func NewBBandsSMA(n int, nStdDev float64) Indicator {
	return &bbands{
		nm:      fmt.Sprintf("BBANDS(SMA, %v)", n),
		n:       int64(n),
		nStdDev: nStdDev,
		isSma:   true,
	}
}

func NewBBandsEMA(n int, nStdDev float64) Indicator {
	return &bbands{
		nm:      fmt.Sprintf("BBANDS(EMA, %v)", n),
		n:       int64(n),
		nStdDev: nStdDev,
		isSma:   false,
	}
}

func (b bbands) name() string {
	return b.nm
}

func (b bbands) yAxisLabel() string {
	return ""
}

func (b bbands) yAxisMin() string {
	return ""
}

func (b bbands) yAxisMax() string {
	return ""
}

func (b bbands) getNumColors() int {
	return 2
}

func (b *bbands) getTitleOpts(top, left int, colorIndex int) []opts.Title {
	b.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[b.ci],
				FontSize: chartLabelFontSize,
			},
			Title: b.nm + "-Ma",
			Left:  px(left),
			Top:   px(top),
		},
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[b.ci+1],
				FontSize: chartLabelFontSize,
			},
			Title: b.nm + "-Upper",
			Left:  px(left),
			Top:   px(top + chartLabelFontHeight),
		},
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[b.ci+1],
				FontSize: chartLabelFontSize,
			},
			Title: b.nm + "-Lower",
			Left:  px(left),
			Top:   px(top + 2*chartLabelFontHeight),
		},
	}
}

func (b bbands) genChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	var u, m, l []float64
	if b.isSma {
		u, m, l = tart.BBandsArr(tart.SMA, closes, b.n, b.nStdDev, b.nStdDev)
	} else {
		u, m, l = tart.BBandsArr(tart.EMA, closes, b.n, b.nStdDev, b.nStdDev)
	}

	uItems := []opts.LineData{}
	mItems := []opts.LineData{}
	lItems := []opts.LineData{}
	for i := 0; i < len(m); i++ {
		uItems = append(uItems, opts.LineData{Value: u[i]})
		mItems = append(mItems, opts.LineData{Value: m[i]})
		lItems = append(lItems, opts.LineData{Value: l[i]})
	}

	ml := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm+"-Ma", mItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[b.ci],
				Opacity: opacityMed,
			}))
	ul := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm+"-Upper", uItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[b.ci+1],
				Opacity: opacityMed,
			}))
	ll := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nm+"-Lower", lItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[b.ci+1],
				Opacity: opacityMed,
			}))

	ml.Overlap(ul, ll)
	return ml
}
