package tachart

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/markcheno/go-talib"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

var (
	ErrInvalidOverlayType = errors.New("unsupported overlay chart type")
)

type IndicatorType string

const (
	SMA  IndicatorType = "SMA"
	EMA  IndicatorType = "EMA"
	MACD IndicatorType = "MACD"
)

type IndicatorConfig struct {
	Type IndicatorType
	// SMA, EMA: "n"
	// MACD: "fast,slow,signal"
	Param string
	// internal
	parsed []int
}

func (c IndicatorConfig) parse() (IndicatorConfig, error) {
	switch c.Type {
	case SMA, EMA:
		n, err := strconv.ParseInt(c.Param, 10, 64)
		if err != nil {
			return IndicatorConfig{}, fmt.Errorf("MA parameter error %v", err)
		}
		c.parsed = append(c.parsed, int(n))
		return c, nil

	case MACD:
		parts := strings.Split(c.Param, ",")
		if len(parts) != 3 {
			return IndicatorConfig{}, fmt.Errorf("unexpected MACD parameters format: fast,slow,signal")
		}
		for _, part := range parts {
			v, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return IndicatorConfig{}, fmt.Errorf("MACD parameter error %v", err)
			}
			c.parsed = append(c.parsed, int(v))
		}
		return c, nil

	default:
		// unknown type
	}

	return IndicatorConfig{}, ErrInvalidOverlayType
}

func (c IndicatorConfig) yAxisLabel() string {
	switch c.Type {
	case SMA, EMA:
		return ""

	case MACD:
		return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
	}

	return ""
}

func (c IndicatorConfig) yAxisMin() string {
	switch c.Type {
	case SMA, EMA:
		return ""

	case MACD:
		return strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", "0", -1)
	}

	return ""
}

func (c IndicatorConfig) yAxisMax() string {
	switch c.Type {
	case SMA, EMA:
		return ""

	case MACD:
		return strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", "0", -1)
	}

	return ""
}

func getChart(vals []float64, xAxis interface{}, c IndicatorConfig, gridIndex int) charts.Overlaper {
	switch c.Type {
	case SMA, EMA:
		var ma []float64
		if c.Type == SMA {
			ma = talib.Sma(vals, c.parsed[0])
		} else {
			ma = talib.Ema(vals, c.parsed[0])
		}
		for i := 0; i < c.parsed[0]; i++ {
			ma[i] = ma[c.parsed[0]]
		}

		items := []opts.LineData{}
		for _, v := range ma {
			items = append(items, opts.LineData{Value: v})
		}

		return charts.NewLine().
			SetXAxis(xAxis).
			AddSeries(fmt.Sprintf("%v%v", c.Type, c.parsed[0]), items, charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}))

	case MACD:
		macd, signal, hist := talib.Macd(vals, c.parsed[0], c.parsed[1], c.parsed[2])

		lineItems := []opts.LineData{}
		for _, v := range macd {
			lineItems = append(lineItems, opts.LineData{Value: v})
		}
		macdLine := charts.NewLine().
			SetXAxis(xAxis).
			AddSeries(fmt.Sprintf("%v(%v)-Diff", c.Type, c.Param), lineItems,
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
			AddSeries(fmt.Sprintf("%v(%v)-Sig", c.Type, c.Param), lineItems,
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
			AddSeries(fmt.Sprintf("%v(%v)", c.Type, c.Param), barItems, charts.WithBarChartOpts(opts.BarChart{
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}))

		macdLine.Overlap(signalLine, histBar)

		return macdLine

	default:
		// will NOT happen
	}

	return nil
}
