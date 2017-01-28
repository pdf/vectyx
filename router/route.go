package router

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gopherjs/vecty"
)

type handlerOrFunc interface{}

type route struct {
	regexp      *regexp.Regexp
	handler     Handler
	handlerFunc HandlerFunc
	score       int
	context     *Context
	lastRender  vecty.ComponentOrHTML
}

func (r *route) match(path string) (score int, context *Context) {
	u, err := url.Parse(path)
	if err != nil {
		return -1, nil
	}
	// Pin regexp to trailing slash
	if u.Path[len(u.Path)-1] != '/' {
		u.Path += `/`
	}
	params := u.Query()

	matches := r.regexp.FindStringSubmatch(u.Path)
	if len(matches) == 0 {
		return -1, nil
	}

	for i, name := range r.regexp.SubexpNames() {
		if name == `` {
			continue
		}
		if _, ok := params[name]; !ok {
			params[name] = make([]string, 0)
		}
		params[name] = append(params[name], matches[i])
	}

	context = &Context{
		Path:   u.Path,
		Params: params,
	}

	return r.score, context
}

func (r *route) render(context *Context) vecty.ComponentOrHTML {
	if r.context != nil && r.context.Children == context.Children && !r.dirtyParams(context.Params) {
		return r.lastRender
	}

	context.Rendered = r.context != nil
	r.context = context

	var result vecty.ComponentOrHTML

	if r.handler != nil {
		r.handler.OnRoute(*r.context)
		result = r.handler
	} else if r.handlerFunc != nil {
		result = r.handlerFunc(*r.context)
	}
	r.lastRender = result

	return result
}

func (r *route) dirtyParams(params url.Values) bool {
	return params.Encode() != r.context.Params.Encode()
}

func compileRegexp(pattern string) *regexp.Regexp {
	u, err := url.Parse(pattern)
	if err != nil {
		panic(err)
	}
	split := strings.Split(u.Path, `/`)
	if len(split) == 1 {
		return regexp.MustCompile(`^/`)
	}
	str := `^`
	for _, s := range split {
		if s == `` {
			continue
		}
		if s[0] == ':' {
			str += fmt.Sprintf("/(?P<%s>[^/]+)", s[1:])
		} else if s == `*` {
			str += `(/.*)?`
		} else {
			str += `/` + s
		}
	}
	// Pin regexp to trailing slash, appended in match()
	str += `/`
	return regexp.MustCompile(str)
}

func newRoute(pattern string, handler handlerOrFunc) *route {
	pattern = strings.TrimSuffix(pattern, `/`)
	regex := compileRegexp(pattern)
	r := &route{
		regexp: regex,
		score:  len(strings.Split(regex.String(), `/`)),
	}
	switch h := handler.(type) {
	case Handler:
		r.handler = h
	case HandlerFunc:
		r.handlerFunc = h
	default:
		panic(`Invalid handler passed to newRoute`)
	}

	return r
}
