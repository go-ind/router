package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/sirupsen/logrus"
)

// ascii art logo
const (
	logoGoind = `
           / ____|/ __ \     |_   _| \ | |  __ \
          | |  __| |  | |______| | |  \| | |  | |
          | | |_ | |  | |______| | |     | |  | |
          | |__| | |__| |     _| |_| |\  | |__| |
           \_____|\____/     |_____|_| \_|_____/ 
           `
)

// Router is a tiny HTTP request router with grouping and middleware support.
type Router struct {
	Routes      map[string]map[string]func(http.ResponseWriter, *http.Request)
	SettingLogs string
	middlewares []func(http.Handler) http.Handler
}

// Group represents a collection of routes sharing the same prefix.
type Group struct {
	router *Router
	prefix string
}

func init() {
	fmt.Println("Minimalist Framework Faster with ")
	fmt.Println(logoGoind)
}

// SetupDefaultRouter returns router with logging enabled.
func SetupDefaultRouter() *Router {
	return &Router{
		Routes:      make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
		SettingLogs: "FULLY_LOGING",
	}
}

// SetupWithNoLogging returns router without logging.
func SetupWithNoLogging() *Router {
	return &Router{
		Routes:      make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
		SettingLogs: "NON_LOGING",
	}
}

// internal function to register route
func (r *Router) addRoute(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, ok := r.Routes[path]; !ok {
		r.Routes[path] = make(map[string]func(http.ResponseWriter, *http.Request))
	}
	r.Routes[path][method] = handler
}

// HTTP method helpers -------------------------------------------------------

func (r *Router) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodGet, path, handler)
}

func (r *Router) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPost, path, handler)
}

func (r *Router) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPut, path, handler)
}

func (r *Router) Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPatch, path, handler)
}

func (r *Router) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodDelete, path, handler)
}

// Group route helpers -------------------------------------------------------

func (g *Group) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	g.router.addRoute(HttpMethodGet, g.prefix+path, handler)
}

func (g *Group) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	g.router.addRoute(HttpMethodPost, g.prefix+path, handler)
}

func (g *Group) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	g.router.addRoute(HttpMethodPut, g.prefix+path, handler)
}

func (g *Group) Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	g.router.addRoute(HttpMethodPatch, g.prefix+path, handler)
}

func (g *Group) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	g.router.addRoute(HttpMethodDelete, g.prefix+path, handler)
}

// Group creates a new route group with the given prefix.
func (r *Router) Group(prefix string) *Group {
	return &Group{router: r, prefix: prefix}
}

// Use appends middlewares to the router. They are executed in the order
// they are added.
func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// ServeHTTP implements http.Handler and dispatches the request to the matching route.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	start := time.Now().In(loc)
	req = StartRecord(req, start)

	handlers, ok := r.Routes[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
		return
	}

	handler, ok := handlers[req.Method]
	if !ok {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h := http.Handler(http.HandlerFunc(handler))
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}
	h.ServeHTTP(w, req)
}

// DumpRequest returns the http request as a string.
func DumpRequest(req *http.Request) string {
	header, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "cannot dump request"
	}
	trim := bytes.ReplaceAll(header, []byte("\r\n"), []byte("   "))
	return string(trim)
}

// JSONMarshal marshals the given type without HTML escaping.
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func Logger(ctx context.Context, response string, statuscode int) {
	var level string

	v, ok := ctx.Value(LogKey).(*Data)
	if ok {
		t := time.Since(v.TimeStart)
		if statuscode >= 200 && statuscode < 400 {
			level = "INFO"
		} else if statuscode >= 400 && statuscode < 500 {
			level = "WARN"
		} else {
			level = "ERROR"
		}

		v.StatusCode = statuscode
		v.Response = response
		v.ExecTime = t.Seconds()

		if statuscode == 0 {
			v.StatusCode = 200
		}

		Output(v, level)
	}
}

// Output prints structured log data.
func Output(out *Data, level string) {
	result, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		fmt.Println("error")
	}
	log.Printf("%+v\n", string(result))
}

// UTCFormatter ...
type UTCFormatter struct {
	logrus.Formatter
}
