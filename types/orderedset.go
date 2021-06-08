package types

import (
	"html/template"
)

// OrderedSet represents an ordered set.
type OrderedSet struct {
	filter map[template.URL]bool
	Values []template.URL
}

// Init creates a new OrderedSet instance, and adds any given items into this set.
func (o *OrderedSet) Init(items ...string) {
	o.filter = make(map[template.URL]bool)
	for _, item := range items {
		o.Add(item)
	}
}

// Add adds a new item into the ordered set.
func (o *OrderedSet) Add(item string) {
	v := template.URL(item)
	if !o.filter[v] {
		o.filter[v] = true
		o.Values = append(o.Values, v)
	}
}
