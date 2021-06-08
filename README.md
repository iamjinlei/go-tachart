<p align="center">
	<img src="https://user-images.githubusercontent.com/19553554/52535979-c0d0e680-2d8f-11e9-85c8-2e9f659e7c6f.png" width=300 height=300 />
</p>

<h1 align="center">go-tachart</h1>

<!--
<p align="center">
    <a href="https://travis-ci.org/go-echarts/go-echarts">
        <img src="https://travis-ci.org/go-echarts/go-echarts.svg?branch=master" alt="Build Status">
    </a>
    <a href="https://goreportcard.com/report/github.com/go-echarts/go-echarts">
        <img src="https://goreportcard.com/badge/github.com/go-echarts/go-echarts" alt="Go Report Card">
    </a>
	<a href="https://github.com/go-echarts/go-echarts/pulls">
        <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="Contributions welcome">
    </a>
    <a href="https://opensource.org/licenses/MIT">
        <img src="https://img.shields.io/badge/License-MIT-brightgreen.svg" alt="MIT License">
    </a>
    <a href="https://pkg.go.dev/github.com/go-echarts/go-echarts/v2">
        <img src="https://godoc.org/github.com/go-echarts/go-echarts?status.svg" alt="GoDoc">
    </a>
</p>
-->

This is a fork and extension to the beautiful [go-echarts](https://github.com/go-echarts/go-echarts) project.
Some of the [go-echarts](https://github.com/go-echarts/go-echarts) code is modified to tailor to the needs of TA charts.
To keep iteration simple, [go-echarts](https://github.com/go-echarts/go-echarts) code is replicated in this repo.
This is still a work-in-progress.
More TA chart and event types will be added to support a wide-range of use cases.

### How It Looks Like

#### Candlestick chart with moving average overlay on top
![Screen Shot 2021-06-08 at 12 05 22 PM](https://user-images.githubusercontent.com/6463139/121121323-f566e700-c851-11eb-9b54-9eb52b0a00d8.png)


#### Candlestick chart with moving average overlay on top + additional indicators
!![Screen Shot 2021-06-08 at 12 04 45 PM](https://user-images.githubusercontent.com/6463139/121121350-ff88e580-c851-11eb-8857-8691c2bb7925.png)


### Usage

In this example, a simple chart is created with a few lines of code. A complete code example can be found in the example folder.

```golang
package main

import (
	"github.com/iamjinlei/go-tachart/tachart"
)

func main() {
	cdls := []tachart.Candle{
		{Label: "2018/1/24", O: 2320.26, C: 2320.26, L: 2287.3, H: 2362.94, V: 149092},
		{Label: "2018/1/25", O: 2300, C: 2291.3, L: 2288.26, H: 2308.38, V: 189092},
		{Label: "2018/1/28", O: 2295.35, C: 2346.5, L: 2295.35, H: 2346.92, V: 159034},
		{Label: "2018/1/29", O: 2347.22, C: 2358.98, L: 2337.35, H: 2363.8, V: 249910},
		{Label: "2018/1/30", O: 2360.75, C: 2382.48, L: 2347.89, H: 2383.76, V: 119910},
		{Label: "2018/1/31", O: 2383.43, C: 2385.42, L: 2371.23, H: 2391.82, V: 89940},
                    ... // To fill more data
		{Label: "2018/6/13", O: 2190.1, C: 2148.35, L: 2126.22, H: 2190.1, V: 239510},
	}

	events := []tachart.Event{
		{
			Type:        tachart.Short,
			Label:       cdls[40].Label,
			Description: "This is a demo event description. Randomly pick this candle to go short on " + cdls[40].Label,
		},
	}

	cfg := tachart.NewConfig().
		SetWidth(1080).
		SetHeight(800).
		AddOverlay(
			tachart.NewSMA(5),
			tachart.NewSMA(20),
		).
		AddIndicator(
			tachart.NewMACD(12, 26, 9),
		).
		UseRepoAssets() // serving assets file from current repo, avoid network access

	c := tachart.New(*cfg)
	c.GenStatic(cdls, events, "/Volumes/tmpfs/tmp/kline.html")
}
```
