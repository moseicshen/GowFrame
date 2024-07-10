package gow

import (
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup                //has all ability of router group
	groups       []*RouterGroup //store all group
}

func New() *Engine {
	group := newRootGroup()
	engine := &Engine{
		RouterGroup: group,
		groups:      []*RouterGroup{group},
	}
	return engine
}

func (engine *Engine) Group(prefix string) *RouterGroup {
	group := engine.group(prefix)
	engine.groups = append(engine.groups, group)
	return group
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.RouterGroup.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
