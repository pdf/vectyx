package router

import (
	"net/url"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

// anchor is a wrapper component for linking content
type anchor struct {
	vecty.Core
	path     string
	params   url.Values
	children vecty.ComponentOrHTML
}

func (a *anchor) onClick(_ *vecty.Event) {
	if a.params != nil {
		GoWithParams(a.path, a.params)
		return
	}
	Go(a.path)
}

// Render implements vecty.Component
func (a *anchor) Render() *vecty.HTML {
	return elem.Anchor(
		prop.Href(`javascript:;`),
		event.Click(a.onClick),
		a.children,
	)
}

// LinkWithParams wraps the provided content in an anchor tag that transitions
// to a new location with URL parameters on click.
func LinkWithParams(path string, params url.Values, content vecty.ComponentOrHTML) vecty.MarkupOrComponentOrHTML {
	return &anchor{path: path, params: params, children: content}
}

// Link wraps the provided content in an anchor tag that transitions to a new
// location on click.
func Link(path string, content vecty.ComponentOrHTML) vecty.MarkupOrComponentOrHTML {
	return LinkWithParams(path, nil, content)
}
