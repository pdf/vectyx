package router_test

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"github.com/pdf/vectyx/router"
)

// MyComponent is a Vecty component that accepts children
type MyComponent struct {
	vecty.Core
	children vecty.ComponentOrHTML
}

// OnRoute implements Handler
func (c *MyComponent) OnRoute(ctx router.Context) {
	// Do something here with ctx.Params, etc if desired

	// Update children and re-render if required
	if ctx.ShouldUpdate {
		c.children = ctx.Children
		vecty.Rerender(c)
	}
}

// Render implements vecty.Component
func (c *MyComponent) Render() *vecty.HTML {
	return elem.Div(
		prop.ID(`MyComponent`),
		c.children,
	)
}

// UserIndex placeholder, implement a real component
type UserIndex struct {
	MyComponent
}

// UserShow placeholder, implement a real component
type UserShow struct {
	MyComponent
}

// UserNetwork placeholder, implement a real component
type UserNetwork struct {
	MyComponent
}

// SearchComponent placeholder, implement a real component
type SearchComponent struct {
	MyComponent
}

func Example() {
	r := router.New(nil)
	r.Group(`/users`, &UserIndex{}, func(r *router.Router) {
		r.Handle(`/:userID`, &UserShow{})
		r.Handle(`/:userID/network`, &UserNetwork{})
	})
	r.Handle(`/search/*`, &SearchComponent{})
	r.Handle(`/`, &MyComponent{})

	vecty.RenderBody(r.Body())
}
