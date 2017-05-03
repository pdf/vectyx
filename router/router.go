// Package router implements a client-side router for the rendering of
// components with the Vecty GopherJS framework.
//
// The router supports nested routes, named parameters (`/users/:userID`) and
// wildcards (`/users/*`).
//
// Currently only hash routing is implemented, HTML 5 history (aka pushState)
// support may follow in a future update.
package router

import (
	"net/url"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
)

// DefaultConfig for convenience.
var DefaultConfig = &Config{}

// HandlerFunc allows the use of in-line functions to produce content for routes.
type HandlerFunc func(Context) vecty.ComponentOrHTML

// Handler is a vecty.Component that implements the OnRoute event receiver.
type Handler interface {
	vecty.Component
	// OnRoute is called when route changes require updating the component.
	OnRoute(Context)
}

// Router is a client-side router that can handle the rendering of nested Vecty
// components.
type Router struct {
	config  *Config
	root    *route
	parent  *Router
	pattern string
	routes  []*route
	routers []*Router
}

// Handle registers the handler for a given pattern.
func (r *Router) Handle(pattern string, h Handler) {
	r.routes = append(r.routes, newRoute(r.canonical(pattern), h))
}

// HandleFunc registers the handler function for a given pattern.
func (r *Router) HandleFunc(pattern string, f HandlerFunc) {
	r.routes = append(r.routes, newRoute(r.canonical(pattern), f))
}

// Group registers the root handler for pattern and a set of nested routes.
func (r *Router) Group(pattern string, h Handler, group func(r *Router)) {
	router := New(r.config)
	router.parent = r
	router.pattern = pattern
	r.routers = append(r.routers, router)
	if h != nil {
		router.root = newRoute(r.canonical(pattern), h)
	}
	group(router)
}

// GroupFunc registers the root handler function for pattern and a set of nested
// routes.
func (r *Router) GroupFunc(pattern string, f HandlerFunc, group func(r *Router)) {
	router := New(r.config)
	router.parent = r
	router.pattern = pattern
	router.root = newRoute(r.canonical(pattern), f)
	r.routers = append(r.routers, router)
	group(router)
}

// Body returns the router result wrapped in a body tag, to be passed to
// vecty.RenderBody()
func (r *Router) Body() vecty.Component {
	r.root = newRoute(r.canonical(r.pattern), &body{})
	return r.start().(vecty.Component)
}

// Render returns the router result
func (r *Router) Render() vecty.ComponentOrHTML {
	return r.start()
}

// canonnical walks up the chain of nested routers to produce the canonical
// pattern for the current level
func (r *Router) canonical(pattern string) string {
	return r.walkPath() + pattern
}

// walkPath recursively builds the pattern
func (r *Router) walkPath() string {
	var path string
	if r.parent != nil {
		path = r.parent.walkPath()
	}
	path += r.pattern
	if path == `/` {
		return ``
	}

	return path
}

// hash obtains the current browser location hash component
func hash() string {
	return js.Global.Get(`location`).Get(`hash`).String()
}

// match walks the router returns the matching components for the current path
func (r *Router) match(path string) (score int, result vecty.ComponentOrHTML, context *Context) {
	var (
		max, matchScore int
		children        vecty.ComponentOrHTML
		childContext    *Context
		child           vecty.ComponentOrHTML
		c               *Context
	)

	// Recursively render child routes
	for _, router := range r.routers {
		matchScore, child, c = router.match(path)
		if matchScore > max {
			max = matchScore
			childContext = c
			children = child
		}
	}
	if children != nil {
		score = max
	}

	// Render local routes
	var (
		winner *route
	)
	for _, route := range r.routes {
		matchScore, c = route.match(path)
		if matchScore > max {
			max = matchScore
			winner = route
			context = c
			if childContext != nil {
				for k, v := range childContext.Params {
					if _, ok := context.Params[k]; !ok {
						context.Params[k] = v
						continue
					}
					context.Params[k] = append(context.Params[k], v...)
				}
			}
		}
	}
	if winner != nil {
		context.Children = children
		context.ShouldUpdate = true
		score = max
		result = winner.render(context)
	}
	// No local result, but matched nested route, so use that instead
	if result == nil && children != nil {
		result = children
	}

	if r.root != nil {
		if matchScore, context = r.root.match(path); matchScore > 0 {
			context.Children = result
			context.ShouldUpdate = true
			return score, r.root.render(context), context
		}
	}

	return score, result, context
}

// update triggers a run of the router
func (r *Router) update() vecty.ComponentOrHTML {
	_, result, _ := r.match(currentPathWithParams())
	return result
}

// start initializes the router, ensures we have a hash component and attaches
// the event listener to trigger updates on hash change
func (r *Router) start() vecty.ComponentOrHTML {
	// Redirect to hash route if missing
	if hash() == `` {
		Go(`/`)
	}
	js.Global.Call(`addEventListener`, `hashchange`, r.update, true)
	return r.update()
}

func currentPathWithParams() string {
	return strings.SplitN(hash(), `#`, 2)[1]
}

// CurrentPath returns the current path component of the location or hash fragment
func CurrentPath() string {
	p := currentPathWithParams()
	u, err := url.Parse(p)
	if err != nil {
		return p
	}
	return u.Path
}

// Go transitions the browser to a new location
func Go(path string) {
	GoWithParams(path, nil)
}

// GoWithParams transitions the browser to a new location, with query params
func GoWithParams(path string, params url.Values) {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	u.RawQuery = params.Encode()
	js.Global.Get(`location`).Set(`hash`, u.RequestURI())
}

// New instantiates a new router.  If the config argument is nil, DefaultConfig
// will be used.
func New(config *Config) *Router {
	if config == nil {
		config = DefaultConfig
	}
	r := &Router{
		config:  config,
		routes:  make([]*route, 0),
		routers: make([]*Router, 0),
	}
	return r
}
