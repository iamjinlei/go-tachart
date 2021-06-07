package tachart

type Config struct {
	w          int
	h          int
	precision  int // decimal places of floating nubmers shown on chart
	overlays   []IndicatorConfig
	indicators []IndicatorConfig
}

func NewConfig() *Config {
	return &Config{
		w:          900,
		h:          500,
		precision:  2,
		overlays:   []IndicatorConfig{},
		indicators: []IndicatorConfig{},
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
