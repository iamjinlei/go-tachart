package tachart

import (
	"errors"
	"fmt"
	"os"

	"github.com/markcheno/go-talib"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/components"
	"github.com/iamjinlei/go-tachart/opts"
)

const (
	tooltipPositionFunc = `
		function (pos, params, el, elRect, size) {
			var obj = {top: 10};
			if (pos[0] > size.viewSize[0]/2) {
				obj['left'] = 30;
			} else {
				obj['right'] = 30;
			}
			return obj;
		}`

	colorDownBar = "#00da3c"
	colorUpBar   = "#ec0000"
)

var (
	ErrDuplicateCandleLabel = errors.New("candles with duplicated labels")
)

type TAChart struct {
	// TODO: support dynamic auto-refresh
	overlays []OverlayChart
}

func New(overlays []OverlayChart) *TAChart {
	return &TAChart{
		overlays: overlays,
	}
}

func (c TAChart) GenStatic(cdls []Candle, events []Event, path string) error {
	x := make([]string, 0)
	klineSeries := []opts.KlineData{}
	volSeries := []opts.BarData{}
	closes := []float64{}
	cdlMap := map[string]*Candle{}
	for _, cdl := range cdls {
		x = append(x, cdl.Label)
		// open,close,low,high
		klineSeries = append(klineSeries, opts.KlineData{Value: []float64{cdl.O, cdl.C, cdl.L, cdl.H}})
		closes = append(closes, cdl.C)

		style := &opts.ItemStyle{
			Color: colorUpBar,
		}
		if cdl.O > cdl.C {
			style = &opts.ItemStyle{
				Color: colorDownBar,
			}
		}
		volSeries = append(volSeries, opts.BarData{
			Value:     cdl.V,
			ItemStyle: style,
		})

		if cdlMap[cdl.Label] != nil {
			return ErrDuplicateCandleLabel
		}
		c := cdl
		cdlMap[cdl.Label] = &c
	}

	chartOpts := []charts.SeriesOpts{
		charts.WithKlineChartOpts(opts.KlineChart{
			BarWidth:   "60%",
			XAxisIndex: 0,
			YAxisIndex: 0,
		}),
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        colorUpBar,
			Color0:       colorDownBar,
			BorderColor:  colorUpBar,
			BorderColor0: colorDownBar,
		}),
	}
	for _, e := range events {
		yPos := e.Position
		if yPos == 0 {
			cdl := cdlMap[e.Label]
			switch e.Type {
			case Long:
				yPos = cdl.L * 0.9
			case Short:
				yPos = cdl.H * 1.01
			default:
				continue
			}
		}
		chartOpts = append(chartOpts, charts.WithMarkPointNameCoordItemOpts(opts.MarkPointNameCoordItem{
			Symbol:     "roundRect",
			SymbolSize: 16,
			Coordinate: []interface{}{e.Label, yPos},
			Label:      eventLabelMap[e.Type].label,
			ItemStyle:  eventLabelMap[e.Type].style,
		}))
	}

	chart := charts.NewKLine().
		SetXAxis(x).
		AddSeries("kline", klineSeries, chartOpts...)

	for _, ol := range c.overlays {
		var vals []float64
		name := ""
		switch ol.Type {
		case SMA:
			vals = talib.Sma(closes, ol.N)
			for i := 0; i < ol.N; i++ {
				vals[i] = vals[ol.N]
			}
			name = fmt.Sprintf("%v%v", ol.Type, ol.N)
		case EMA:
			vals = talib.Ema(closes, ol.N)
			for i := 0; i < ol.N; i++ {
				vals[i] = vals[ol.N]
			}
			name = fmt.Sprintf("%v%v", ol.Type, ol.N)
		default:
			continue
		}

		items := []opts.LineData{}
		for _, v := range vals {
			items = append(items, opts.LineData{Value: v})
		}

		line := charts.NewLine().
			SetXAxis(x).
			AddSeries(name, items, charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: 0,
				YAxisIndex: 0,
				ZLevel:     100,
			}))
		chart.Overlap(line)
	}

	bar := charts.NewBar().
		SetXAxis(x).
		AddSeries("vol", volSeries, charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: 1,
			YAxisIndex: 1,
		}))
	chart.Overlap(bar)

	chart.ExtendXAxis(opts.XAxis{
		GridIndex:   1,
		SplitNumber: 20,
		Data:        x,
		AxisTick: &opts.AxisTick{
			Show: false,
		},
		AxisLabel: &opts.AxisLabel{
			Show: false,
		},
	})
	chart.ExtendYAxis(opts.YAxis{
		GridIndex:   1,
		Scale:       true,
		SplitNumber: 2,
		SplitLine: &opts.SplitLine{
			Show: false,
		},
		AxisLabel: &opts.AxisLabel{
			Show:         true,
			ShowMaxLabel: true,
		},
		Min: 0,
		Max: "dataMax",
	})
	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "candles",
		}),
		charts.WithAxisPointerOpts(opts.AxisPointer{
			Type: "line",
			Snap: true,
			Link: opts.AxisPointerLink{
				XAxisIndex: "all",
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			GridIndex:   0,
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			GridIndex: 0,
			Scale:     true,
			SplitArea: &opts.SplitArea{
				Show: true,
			},
		}),
		charts.WithGridOpts(opts.Grid{
			Left:   "10%",
			Right:  "8%",
			Height: "50%",
		},
			opts.Grid{
				Left:   "10%",
				Right:  "8%",
				Top:    "73%",
				Height: "16%",
			}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0, 1},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "mousemove|click",
			Position:  opts.FuncOpts(tooltipPositionFunc),
		}),
	)

	page := components.NewPage().AddCharts(chart)
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	return page.Render(fp)
}
