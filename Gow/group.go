package gow

import "net/http"

type RouterGroup struct {
	*router
	prefix     string
	middleware []HandlerFunc
}

func newRootGroup() *RouterGroup {
	return &RouterGroup{
		router: newRouter(),
		prefix: "",
	}
}

func (group *RouterGroup) group(prefix string) *RouterGroup {
	return &RouterGroup{
		router: group.router,
		prefix: group.prefix + prefix,
	}
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middleware = append(group.middleware, middlewares...)
}

func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	group.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) {
	//absolutePath := group.prefix + relativePath
}
