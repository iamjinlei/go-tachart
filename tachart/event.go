package tachart

import (
	"fmt"

	"github.com/iamjinlei/go-tachart/opts"
)

type EventType byte

const (
	Long  EventType = 'L'
	Short EventType = 'S'
	Open  EventType = 'O'
	Close EventType = 'C'
)

type eventStyle struct {
	label *opts.Label
	style *opts.ItemStyle
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
		},
	}
)

type Event struct {
	Type        EventType
	Label       string // x-axis label. Should match to one of the candles
	Description string // any user-defined description wants to appear on tooltip
}
