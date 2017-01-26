package router

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

// body is a wrapper component
type body struct {
	vecty.Core
	children vecty.ComponentOrHTML
}

// OnRoute implements Handler
func (b *body) OnRoute(ctx Context) {
	b.children = ctx.Children
	if ctx.Rendered {
		vecty.Rerender(b)
	}
}

// Render implements vecty.Component
func (b *body) Render() *vecty.HTML {
	return elem.Body(
		prop.ID(`Body`),
		b.children,
	)
}
