package components

import (
	"github.com/iamjinlei/go-tachart/opts"
	"github.com/iamjinlei/go-tachart/render"
)

// Charter
type Charter interface {
	Type() string
	GetAssets() opts.Assets
	Validate()
}

// Page represents a page chart.
type Page struct {
	render.Renderer
	opts.Initialization
	opts.Assets

	PageBackgroundColor string
	Layout              Layout
	Charts              []interface{}
	ChartArea           ChartLayout
}

// NewPage creates a new page.
func NewPage(assetsHost string) *Page {
	page := &Page{}
	page.Assets.InitAssets()

	page.Renderer = render.NewPageRender(page, page.Validate)
	page.AssetsHost = assetsHost
	page.ChartArea = PageCenterLayout
	page.PageBackgroundColor = "#FFFFFF"

	return page
}

func (page *Page) SetLayout(layout Layout) *Page {
	page.Layout = layout
	return page
}

func (page *Page) SetBackgroundColor(color string) *Page {
	page.PageBackgroundColor = color
	return page
}

// AddCharts adds new charts to the page.
func (page *Page) AddCharts(charts ...Charter) *Page {
	for i := 0; i < len(charts); i++ {
		assets := charts[i].GetAssets()
		for _, v := range assets.JSAssets.Values {
			page.JSAssets.Add(string(v))
		}

		for _, v := range assets.CSSAssets.Values {
			page.CSSAssets.Add(string(v))
		}
		charts[i].Validate()
		page.Charts = append(page.Charts, charts[i])
	}
	return page
}

// Validate
func (page *Page) Validate() {
	page.Initialization.Validate()
	page.Assets.Validate(page.AssetsHost)
}
