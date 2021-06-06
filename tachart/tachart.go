package tachart

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/markcheno/go-talib"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/components"
	"github.com/iamjinlei/go-tachart/opts"
)

const (
	tooltipPositionFunc = `
		function(pos, params, el, elRect, size) {
			var obj = {top: 10};
			if (pos[0] > size.viewSize[0]/2) {
				obj['left'] = 30;
			} else {
				obj['right'] = 30;
			}
			return obj;
		}`
	tooltipFormatterFuncTpl = `
		function(value) {
			var eventMap = JSON.parse('__EVENT_MAP__');

			var title = (sz,txt) => '<span style="display:inline;line-height:'+(sz+2)+'px;font-size:'+sz+'px;font-weight:bold;">'+txt+'</span>';
			var square = (sz,sign,color,txt) => '<span style="display:inline;line-height:'+(sz+2)+'px;font-size:'+sz+'px;"><span style="display:inline-block;height:'+(sz+2)+'px;border-radius:3px;padding:1px 4px 1px 4px;text-align:center;margin-right:10px;background-color:' + color + ';vertical-align:top;">'+sign+'</span>'+txt+'</span>';
			var wrap = (sz,txt,width) => '<span style="display:inline-block;width:'+width+'px;word-break:break-word;word-wrap:break-word;white-space:pre-wrap;line-height:'+(sz+2)+'px;font-size:'+sz+'px;">'+txt+'</span>';

			value.sort((a, b) => a.seriesIndex -b.seriesIndex);
			var cdl = value[0];
			var ret = title(14,cdl.axisValueLabel) + '<br/>' +
			square(13,'O',cdl.color,cdl.value[1].toFixed(__DECIMAL_PLACES__)) + '<br/>' +
			square(13,'C',cdl.color,cdl.value[2].toFixed(__DECIMAL_PLACES__)) + '<br/>' +
			square(13,'L',cdl.color,cdl.value[3].toFixed(__DECIMAL_PLACES__)) + '<br/>' +
			square(13,'H',cdl.color,cdl.value[4].toFixed(__DECIMAL_PLACES__)) + '<br/>';
			for (var i = 1; i < value.length; i++) {
				var s = value[i];
				ret += square(13,s.seriesName,s.color,s.value.toFixed(__DECIMAL_PLACES__)) + '<br/>';
			}

			var desc = eventMap[cdl.axisValueLabel];
			if (desc) {
				ret += '<br/>' + title(14,'Event Desc.') + '<br/>' + wrap(13,desc,160);
			}
			return ret;
		}`
	candleYMinFuncTpl = `
		function(value) {
			var m = Math.pow(10, __DECIMAL_PLACES__);
			return Math.round(value.min*0.98 * m) / m;
		}`
	yLabelFormatterFuncTpl = `
		function(value) {
			return value.toFixed(__DECIMAL_PLACES__);
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

func (c TAChart) GenStatic(title string, cdls []Candle, events []Event, path string) error {
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

	chart := charts.NewKLine().SetXAxis(x).AddSeries("kline",
		klineSeries,
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
		}))

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

	decimalPlaces := fmt.Sprintf("%v", maxDecimalPlaces(cdls))
	candleYMinFunc := strings.Replace(candleYMinFuncTpl, "__DECIMAL_PLACES__", decimalPlaces, -1)
	yLabelFormatterFunc := strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", decimalPlaces, -1)
	tooltipFormatterFunc := strings.Replace(tooltipFormatterFuncTpl, "__DECIMAL_PLACES__", decimalPlaces, -1)
	eventDescMap := map[string]string{}
	for _, e := range events {
		eventDescMap[e.Label] = e.Description
	}
	fmt.Printf("%v\n", toJson(eventDescMap))
	tooltipFormatterFunc = strings.Replace(tooltipFormatterFunc, "__EVENT_MAP__", toJson(eventDescMap), 1)
	fmt.Printf("%v\n", tooltipFormatterFunc)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithAxisPointerOpts(opts.AxisPointer{
			Type: "line",
			Snap: true,
			Link: opts.AxisPointerLink{
				XAxisIndex: "all",
			},
		}),
		charts.WithGridOpts(
			opts.Grid{
				Left:   "10%",
				Right:  "8%",
				Height: "55%",
			},
			opts.Grid{
				Left:   "10%",
				Right:  "8%",
				Top:    "60%",
				Height: "3%",
			},
			opts.Grid{
				Left:   "10%",
				Right:  "8%",
				Top:    "75%",
				Height: "13%",
			},
		),
		charts.WithXAxisOpts(opts.XAxis{
			Show:        true,
			GridIndex:   0,
			SplitNumber: 20,
		}, 0),
		charts.WithYAxisOpts(opts.YAxis{
			Show:      true,
			GridIndex: 0,
			Scale:     true,
			SplitArea: &opts.SplitArea{
				Show: true,
			},
			Min: opts.FuncOpts(candleYMinFunc),
			Max: "dataMax",
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				ShowMinLabel: true,
				ShowMaxLabel: true,
				Formatter:    opts.FuncOpts(yLabelFormatterFunc),
			},
		}, 0),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0, 1, 2},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "mousemove|click",
			Position:  opts.FuncOpts(tooltipPositionFunc),
			Formatter: opts.FuncOpts(tooltipFormatterFunc),
		}),
	)
	chart.ExtendXAxis(
		opts.XAxis{
			Show:      false,
			GridIndex: 1,
			//SplitNumber: 20,
			Data: x,
			/*
				AxisTick: &opts.AxisTick{
					Show: false,
				},
				AxisLabel: &opts.AxisLabel{
					Show: false,
				},
			*/
		},
		opts.XAxis{
			Show:        true,
			GridIndex:   2,
			SplitNumber: 20,
			Data:        x,
			AxisTick: &opts.AxisTick{
				Show: false,
			},
			AxisLabel: &opts.AxisLabel{
				Show: false,
			},
		})
	chart.ExtendYAxis(
		opts.YAxis{
			Show:      false,
			GridIndex: 1,
			/*
				SplitLine: &opts.SplitLine{
					Show: false,
				},
				AxisLabel: &opts.AxisLabel{
					Show: false,
				},
			*/
		},
		opts.YAxis{
			Show:        true,
			GridIndex:   2,
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

	/*
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
	*/
	chartOpts := []charts.SeriesOpts{
		charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: 1,
			YAxisIndex: 1,
		}),
	}
	for _, e := range events {
		chartOpts = append(chartOpts, charts.WithMarkPointNameCoordItemOpts(opts.MarkPointNameCoordItem{
			Symbol:     "roundRect",
			SymbolSize: 16,
			Coordinate: []interface{}{e.Label, 0},
			Label:      eventLabelMap[e.Type].label,
			ItemStyle:  eventLabelMap[e.Type].style,
		}))
	}
	event := charts.NewBar().AddSeries("events", []opts.BarData{}, chartOpts...)
	chart.Overlap(event)

	bar := charts.NewBar().
		SetXAxis(x).
		AddSeries("Vol", volSeries, charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: 2,
			YAxisIndex: 2,
		}))
	chart.Overlap(bar)

	page := components.NewPage().AddCharts(chart)
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	return page.Render(fp)
}

func toJson(o interface{}) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(o)
	return string(buf.Bytes())
}
