package router

import (
	"net/url"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

// anchor is a wrapper component for linking content
type anchor struct {
	vecty.Core
	path     string
	params   url.Values
	children vecty.ComponentOrHTML
}

func (a *anchor) onClick(_ vecty.Event) {
	Go(a.path, a.params)
}

// Render implements vecty.Component
func (a *anchor) Render() *vecty.HTML {
	return elem.Anchor(
		prop.Href(`javascript:;`),
		a.children,
	)
}
