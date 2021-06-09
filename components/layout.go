package components

import (
	"html/template"
)

type ChartLayout string

const (
	PageNoneLayout   ChartLayout = "none"
	PageCenterLayout ChartLayout = "center"
	PageFlexLayout   ChartLayout = "flex"
)

type Layout struct {
	TemplateColumns template.CSS
	TopHeight       template.CSS
	BottomHeight    template.CSS
	TopContent      template.HTML
	BottomContent   template.HTML
	LeftContent     template.HTML
	RightContent    template.HTML
}
