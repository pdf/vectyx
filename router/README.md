# [Vecty](https://github.com/gopherjs/vecty) Router [![GoDoc](https://godoc.org/github.com/pdf/vectyx/router?status.svg)](http://godoc.org/github.com/pdf/vectyx/router) ![License-MIT](http://img.shields.io/badge/license-MIT-red.svg)

# router
--
    import "github.com/pdf/vectyx/router"

Package router implements a client-side router for the rendering of components
with the Vecty GopherJS framework.

The router supports nested routes, named parameters (`/users/:userID`) and
wildcards (`/users/*`).

Currently only hash routing is implemented, HTML 5 history (aka pushState)
support may follow in a future update.

## Usage

```go
var DefaultConfig = &Config{}
```
DefaultConfig for convenience.

#### func  Go

```go
func Go(path string, params url.Values)
```
Go transitions the browser to a new location

#### func  Link

```go
func Link(path string, params url.Values, content vecty.ComponentOrHTML) vecty.ComponentOrHTML
```
Link wraps the provided content in an anchor tag that transitions to a new
location on click

#### type Config

```go
type Config struct {
}
```

Config provides configuration state for the Router

TODO: Currently there are no configuration options.

#### type Context

```go
type Context struct {
	// Path is the currently matching URL component.
	Path string
	// Params may be populated via named router components, or query params.
	Params url.Values
	// Children optionally contains any nested components that should be rendered
	// by the handler.
	Children vecty.ComponentOrHTML
	// Rendered is `false` if this is the first time the component has been added
	// to the tree.  If Rendered is `true`, the component should be re-rendered
	// (ie - `vecty.Rerender()`) if an update is desired.
	Rendered bool
}
```

Context is provided to handlers for action when routes change.

#### type Handler

```go
type Handler interface {
	vecty.Component
	// OnRoute is called when route changes require updating the component.
	OnRoute(Context)
}
```

Handler is a vecty.Component that implements the OnRoute event receiver.

#### type HandlerFunc

```go
type HandlerFunc func(Context) vecty.ComponentOrHTML
```

HandlerFunc allows the use of in-line functions to produce content for routes.

#### type Router

```go
type Router struct {
}
```

Router is a client-side router that can handle the rendering of nested Vecty
components.

#### func  New

```go
func New(config *Config) *Router
```
New instantiates a new router. If the config argument is nil, DefaultConfig will
be used.

#### func (*Router) Body

```go
func (r *Router) Body() vecty.Component
```
Body returns the router result wrapped in a body tag, to be passed to
vecty.RenderBody()

#### func (*Router) Group

```go
func (r *Router) Group(pattern string, h Handler, group func(r *Router))
```
Group registers the root handler for pattern and a set of nested routes.

#### func (*Router) GroupFunc

```go
func (r *Router) GroupFunc(pattern string, f HandlerFunc, group func(r *Router))
```
GroupFunc registers the root handler function for pattern and a set of nested
routes.

#### func (*Router) Handle

```go
func (r *Router) Handle(pattern string, h Handler)
```
Handle registers the handler for a given pattern.

#### func (*Router) HandleFunc

```go
func (r *Router) HandleFunc(pattern string, f HandlerFunc)
```
HandleFunc registers the handler function for a given pattern.

#### func (*Router) Render

```go
func (r *Router) Render() vecty.ComponentOrHTML
```
Render returns the router result

## Example
```go
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
	c.children = ctx.Children
	if ctx.Rendered {
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

func main() {
	r := router.New(nil)
	r.Group(`/users`, &UserIndex{}, func(r *router.Router) {
		r.Handle(`/:userID`, &UserShow{})
		r.Handle(`/:userID/network`, &UserNetwork{})
	})
	r.Handle(`/search/*`, &SearchComponent{})
	r.Handle(`/`, &MyComponent{})

	vecty.RenderBody(r.Body())
}
```
