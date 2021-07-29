package tachart

import (
	"fmt"
	"strings"

	"github.com/iamjinlei/go-tachart/charts"
	"github.com/iamjinlei/go-tachart/opts"
)

type line struct {
	nms     []string
	valsArr [][]float64
	nc      int
	ci      int
	dp      int
}

func NewLine(name string, vals []float64) Indicator {
	return &line{
		nms:     []string{name},
		valsArr: [][]float64{vals},
		nc:      1,
		dp:      decimals(vals),
	}
}

func NewLine2(n0 string, vals0 []float64, n1 string, vals1 []float64) Indicator {
	return &line{
		nms:     []string{n0, n1},
		valsArr: [][]float64{vals0, vals1},
		nc:      2,
		dp:      decimals(vals0, vals1),
	}
}

func NewLine3(n0 string, vals0 []float64, n1 string, vals1 []float64, n2 string, vals2 []float64) Indicator {
	return &line{
		nms:     []string{n0, n1, n2},
		valsArr: [][]float64{vals0, vals1, vals2},
		nc:      3,
		dp:      decimals(vals0, vals1, vals2),
	}
}

func (b line) name() string {
	return strings.Join(b.nms, ", ")
}

func (b line) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b line) yAxisMin() string {
	return strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b line) yAxisMax() string {
	return strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", fmt.Sprintf("%v", b.dp), -1)
}

func (b line) getNumColors() int {
	return b.nc
}

func (b *line) getTitleOpts(top, left int, colorIndex int) []opts.Title {
	b.ci = colorIndex
	var tls []opts.Title
	for i, nm := range b.nms {
		tls = append(tls, opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[b.ci+i],
				FontSize: chartLabelFontSize,
			},
			Title: nm,
			Left:  px(left),
			Top:   px(top + i*chartLabelFontHeight),
		})
	}
	return tls
}

func (b line) genChart(_, _, _, _, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {
	lineItems := []opts.LineData{}
	for _, v := range b.valsArr[0] {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}

	c := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(b.nms[0], lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[b.ci],
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

	for i := 1; i < len(b.nms); i++ {
		lineItems := []opts.LineData{}
		for _, v := range b.valsArr[i] {
			lineItems = append(lineItems, opts.LineData{Value: v})
		}

		line := charts.NewLine().
			SetXAxis(xAxis).
			AddSeries(b.nms[i], lineItems,
				charts.WithLineChartOpts(opts.LineChart{
					Symbol:     "none",
					XAxisIndex: gridIndex,
					YAxisIndex: gridIndex,
					ZLevel:     100,
				}),
				charts.WithLineStyleOpts(opts.LineStyle{
					Color:   colors[b.ci+i],
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
		c.Overlap(line)
	}

	return c
}
