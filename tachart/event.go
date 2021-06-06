package tachart

import (
	"fmt"

	"github.com/iamjinlei/go-tachart/opts"
)

type EventType byte

const (
	Long  EventType = 'L'
	Short EventType = 'S'
)

type eventStyle struct {
	label *opts.Label
	style *opts.ItemStyle
}

var (
	eventLabelMap = map[EventType]*eventStyle{
		Long: &eventStyle{
			label: &opts.Label{
				Show:  true,
				Color: "#FFFFFF",
				//BackgroundColor: colorUpBar,
				Formatter: fmt.Sprintf("%c", Long),
			},
			style: &opts.ItemStyle{
				Color:       colorUpBar,
				BorderColor: colorUpBar,
			},
		},
		Short: &eventStyle{
			label: &opts.Label{
				Show:  true,
				Color: "#0000FF",
				// BackgroundColor: colorDownBar,
				Formatter: fmt.Sprintf("%c", Short),
			},
			style: &opts.ItemStyle{
				Color:       colorDownBar,
				BorderColor: colorDownBar,
			},
		},
	}
)

type Event struct {
	Type        EventType
	Label       string // x-axis label. Should match to one of the candles
	Description string // any user-defined description wants to appear on tooltip
}
