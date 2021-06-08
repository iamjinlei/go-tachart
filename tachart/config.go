package tachart

import (
	"path/filepath"
	"runtime"
)

type Config struct {
	w          int
	h          int
	precision  int // decimal places of floating nubmers shown on chart
	overlays   []IndicatorConfig
	indicators []IndicatorConfig
	assetsHost string
}

func NewConfig() *Config {
	return &Config{
		w:          900,
		h:          500,
		precision:  2,
		overlays:   []IndicatorConfig{},
		indicators: []IndicatorConfig{},
		assetsHost: "https://go-echarts.github.io/go-echarts-assets/assets/",
	}
}

func (c *Config) SetWidth(w int) *Config {
	c.w = w
	return c
}

func (c *Config) SetHeight(h int) *Config {
	c.h = h
	return c
}

func (c *Config) SetPrecision(p int) *Config {
	c.precision = p
	return c
}

func (c *Config) AddOverlay(cfgs ...IndicatorConfig) *Config {
	c.overlays = append(c.overlays, cfgs...)
	return c
}

func (c *Config) AddIndicator(cfgs ...IndicatorConfig) *Config {
	c.indicators = append(c.indicators, cfgs...)
	return c
}

func (c *Config) UseRepoAssets() *Config {
	// serving assets from "this" repo in local file system
	// with accessing network
	_, path, _, _ := runtime.Caller(0)
	path = filepath.Dir(path)
	c.assetsHost = filepath.Join("file:/"+filepath.Dir(path), "assets/")
	return c
}

func (c *Config) SetAssetsHost(assetsHost string) *Config {
	// serving assets from specified host
	c.assetsHost = assetsHost
	return c
}
