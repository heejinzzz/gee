package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type H map[string]interface{}

type Context struct {
	ResWriter http.ResponseWriter
	Req       *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middlewares and handlers
	handlersChain []HandlerFunc
	index         int
	isAbort       bool
	// engine pointer
	engine *Engine
}

func newContext(resWriter http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		ResWriter: resWriter,
		Req:       req,
		Path:      req.URL.Path,
		Method:    req.Method,
		index:     -1,
	}
}

// Next process middlewares and handlers behind the index in handlersChain
func (c *Context) Next() {
	c.index++
	n := len(c.handlersChain)
	for ; c.index < n && !c.isAbort; c.index++ {
		handler := c.handlersChain[c.index]
		if handler != nil {
			handler(c)
		}
	}
}

// Abort end processing handlersChain
func (c *Context) Abort() {
	c.isAbort = true
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Status(statusCode int) {
	c.StatusCode = statusCode
	c.ResWriter.WriteHeader(statusCode)
}

func (c *Context) SetHeader(key string, value string) {
	c.ResWriter.Header().Set(key, value)
}

func (c *Context) String(statusCode int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(statusCode)
	c.ResWriter.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(statusCode int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(statusCode)
	encoder := json.NewEncoder(c.ResWriter)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.ResWriter, err.Error(), 500)
	}
}

func (c *Context) Data(statusCode int, data []byte) {
	c.Status(statusCode)
	c.ResWriter.Write(data)
}

func (c *Context) HTML(statusCode int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(statusCode)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.ResWriter, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

func (c *Context) Fail(statusCode int, info string) {
	c.Status(statusCode)
	c.String(statusCode, "%d ERROR: %s", statusCode, c.Req.RequestURI)
	log.Printf("%d ERROR: %s. Error: %s\n", statusCode, c.Req.RequestURI, info)
	c.Abort()
}

func (c *Context) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.ResWriter, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return "", err
	}
	value, _ := url.QueryUnescape(cookie.Value)
	return value, nil
}
