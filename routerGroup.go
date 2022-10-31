package gee

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine // all groups share the same Engine instance
}

// Group is defined to create a new RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, pattern string, handlers []HandlerFunc) {
	engine := group.engine
	engine.addRoute(method, group.prefix+pattern, handlers)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handlers ...HandlerFunc) {
	group.addRoute("GET", pattern, handlers)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handlers ...HandlerFunc) {
	group.addRoute("POST", pattern, handlers)
}

// Use is defined to add middlewares to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Fail(http.StatusNotFound, err.Error())
			return
		}
		c.Status(http.StatusOK)
		fileServer.ServeHTTP(c.ResWriter, c.Req)
	}
}

// Static register static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}
