package gee

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup // store all router groups
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

// New creates an empty engine with no middlewares
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// NewDefault creates an engine with Logger and Recovery as middlewares
func NewDefault() *Engine {
	engine := New()
	engine.Use(Logger, Recovery)
	return engine
}

func (engine *Engine) addRoute(method string, pattern string, handlers []HandlerFunc) {
	engine.router.addRoute(method, pattern, handlers)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handlers ...HandlerFunc) {
	engine.addRoute("GET", pattern, handlers)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handlers ...HandlerFunc) {
	engine.addRoute("POST", pattern, handlers)
}

func (engine *Engine) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
	c := newContext(resWriter, req)
	c.engine = engine

	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c.handlersChain = middlewares

	engine.router.handle(c)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) error {
	log.Printf("[GEE] Engine Start. Start Listening %s\n", addr)
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
