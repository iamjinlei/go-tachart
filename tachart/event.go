package tachart

import (
	"fmt"

	"github.com/iamjinlei/go-tachart/opts"
)

const (
	symbolSize = 16
)

type EventType byte

const (
	Long        EventType = 'L'
	Short       EventType = 'S'
	Open        EventType = 'O'
	Close       EventType = 'C'
	CustomEvent EventType = 'X'
)

type eventStyle struct {
	label      *opts.Label
	style      *opts.ItemStyle
	symbolSize int
}

var (
	eventLabelMap = map[EventType]*eventStyle{
		Long: &eventStyle{
			label: &opts.Label{
				Show:      true,
				Color:     "#FFFFFF",
				Formatter: fmt.Sprintf("%c", Long),
			},
			style: &opts.ItemStyle{
				Color:       colorUpBar,
				BorderColor: colorUpBar,
			},
			symbolSize: symbolSize,
		},
		Short: &eventStyle{
			label: &opts.Label{
				Show:      true,
				Color:     "#0000FF",
				Formatter: fmt.Sprintf("%c", Short),
			},
			style: &opts.ItemStyle{
				Color:       colorDownBar,
				BorderColor: colorDownBar,
			},
			symbolSize: symbolSize,
		},
		Open: &eventStyle{
			label: &opts.Label{
				Show:      true,
				Color:     "#FFFFFF",
				Formatter: fmt.Sprintf("%c", Open),
			},
			style: &opts.ItemStyle{
				Color:       "#1D8348",
				BorderColor: "#1D8348",
			},
			symbolSize: symbolSize,
		},
		Close: &eventStyle{
			label: &opts.Label{
				Show:      true,
				Color:     "#FFFFFF",
				Formatter: fmt.Sprintf("%c", Close),
			},
			style: &opts.ItemStyle{
				Color:       "#943126",
				BorderColor: "#943126",
			},
			symbolSize: symbolSize,
		},
	}
)

type EventMark struct {
	Name        string // mark label string
	FontColor   string // mark label font color
	BgColor     string // mark icon color
	BorderColor string // mark icon border color, default to BgColor if empty
	SymbolSize  int    // symbol size, use default if 0
}

func (m EventMark) toEventStyle() *eventStyle {
	fc := m.FontColor
	bgc := m.BgColor
	bc := m.BorderColor
	if bc == "" {
		bc = bgc
	}
	sz := m.SymbolSize
	if sz == 0 {
		sz = symbolSize
	}

	return &eventStyle{
		label: &opts.Label{
			Show:      true,
			Color:     fc,
			Formatter: m.Name,
		},
		style: &opts.ItemStyle{
			Color:       bgc,
			BorderColor: bc,
		},
		symbolSize: sz,
	}
}

type Event struct {
	Type        EventType
	Label       string // x-axis label. Should match to one of the candles
	Description string // any user-defined description wants to appear on tooltip
	EventMark   EventMark
}
