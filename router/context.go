package router

import (
	"net/url"

	"github.com/gopherjs/vecty"
)

// Context is provided to handlers for action when routes change.
type Context struct {
	// Path is the currently matching URL component.
	Path string
	// Params may be populated via named router components, or query params.
	Params url.Values
	// Children optionally contains any nested components that should be rendered
	// by the handler.
	Children vecty.ComponentOrHTML
	// ShouldUpdate is `true` if the route change affects the component, but always
	// `false` the first time a component is added to the tree. When ShouldUpdate is
	// `true`, the component should be re-rendered (ie - `vecty.Rerender()`) if an
	// update is desired.
	ShouldUpdate bool
}
