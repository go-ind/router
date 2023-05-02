package goindrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// http://patorjk.com/software/taag/#p=display&f=Big&t=GO-IND
	logoGoind = `
	   / ____|/ __ \     |_   _| \ | |  __ \
	  | |  __| |  | |______| | |  \| | |  | |
	  | | |_ | |  | |______| | |     | |  | |
	  | |__| | |__| |     _| |_| |\  | |__| |
	   \_____|\____/     |_____|_| \_|_____/
	   `
)

type (
	group struct {
		Router
		prefix string
		RoutesInterface
		// handlers []func(http.HandlerFunc) http.HandlerFunc
	}
	Router struct {
		Routes      map[string]map[string]func(http.ResponseWriter, *http.Request)
		SettingLogs string
		middlewares []func(http.Handler) http.Handler
		patern      string
		RoutesInterface
		http.Handler
		// Routes
	}
)

func SetupDefaultRouter() *Router {
	r := &Router{
		Routes:      make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
		SettingLogs: "FULLY_LOGING",
	}
	return r
}

func SetupWithNoLogging() *Router {
	r := &Router{
		Routes:      make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
		SettingLogs: "NON_LOGING",
	}
	return r
}

type RoutesInterface interface {
	Get(path string, handler func(http.ResponseWriter, *http.Request))
	Post(path string, handler func(http.ResponseWriter, *http.Request))

	Put(path string, handler func(http.ResponseWriter, *http.Request))
	Patch(path string, handler func(http.ResponseWriter, *http.Request))
	Delete(path string, handler func(http.ResponseWriter, *http.Request))

	Use(handler func(http.ResponseWriter, *http.Request))
}

func init() {
	fmt.Println("Minimalist Framework Faster with ")
	fmt.Println(logoGoind)
}
func (r *Router) addRoute(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, ok := r.Routes[path]; !ok {
		r.Routes[path] = make(map[string]func(http.ResponseWriter, *http.Request))
	}
	r.Routes[path][method] = handler
}
func (r *Router) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	pathFinal := r.patern + path
	r.addRoute(HttpMethodGet, pathFinal, handler)
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

// ////////////////////
func (r *group) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	pathFinal := r.prefix + path
	r.Router.addRoute(HttpMethodGet, pathFinal, handler)
	// r.addRoute(HttpMethodGet, pathFinal, handler)
}
func (r *group) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPost, path, handler)
}

func (r *group) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPut, path, handler)
}

func (r *group) Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPatch, path, handler)
}

func (r *group) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodDelete, path, handler)
}

func (r *Router) Use(midleware ...func(http.Handler) http.Handler) {
	// r.addRoute(HttpMethodDelete, path, handler)
	// if mx.handler != nil {
	// 	panic("chi: all middlewares must be defined before routes on a mux")
	// }
	r.middlewares = append(r.middlewares, midleware...)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	start := time.Now().In(loc)
	req = StartRecord(req, start)

	if handlers, ok := r.Routes[req.URL.Path]; ok {
		if handler, ok := handlers[req.Method]; ok {

			for _, v := range r.middlewares {
				v(r).ServeHTTP(w, req)

			}
			handler(w, req)
			return
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
	http.NotFound(w, req)
}

func DumpRequest(req *http.Request) string {
	header, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "cannot dump request"
	}

	trim := bytes.ReplaceAll(header, []byte("\r\n"), []byte("   "))
	return string(trim)
}

// JSONMarshal is func
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	fmt.Println("Err", err)
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

// Output for output to terminal
func Output(out *Data, level string) {
	result, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		fmt.Println("error")
	}
	fmt.Printf("%+v\n", string(result))
	// logrus.SetFormatter(UTCFormatter{&logrus.JSONFormatter{}})
	// if level == "ERROR" {
	// 	logrus.WithField("data", out).Error("apps")
	// } else if level == "INFO" {
	// 	logrus.WithField("data", out).Info("apps")
	// } else if level == "WARN" {
	// 	logrus.WithField("data", out).Warn("apps")
	// }
}

// UTCFormatter ...
type UTCFormatter struct {
	logrus.Formatter
}

func (r *Router) Group(prefix string, middleware ...http.HandlerFunc) *group {
	return &group{
		Router: *r,
		prefix: prefix,
		// handlers: middleware,
	}
	// r.patern = prefix
	// fmt.Printf()
	// return rs
	// fn(r)
}

// func MultipleMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {

// 	if len(m) < 1 {
// 		return h
// 	}

// 	wrapped := h

//		// loop in reverse to preserve middleware order
//		for i := len(m) - 1; i >= 0; i-- {
//			wrapped = m[i](wrapped)
//		}
//		return wrapped
//	}
func main() {
	r := &Router{
		Routes:      make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
		SettingLogs: "FULLY_LOGING",
	}

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Hello, world!")
	})

	v1 := r.Group("/parent")
	v1.Get("/ase", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Hello, world!")
	})

	v2 := r.Group("/parent2")
	v2.Get("/parent2-child1", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Hello, worldssss!")
		ResponseJSON(w, ctx, 200, true, "Succes", nil, nil)
	})

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		// id := req.URL.Query().Get("id")
		// fmt.Println("id =>", id)
		// fmt.Fprint(w, "POST succes cok")
		ResponseJSON(w, ctx, 200, true, "Succes", nil, nil)
	})

	r.Put("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "PUT succes cok")
	})

	r.Delete("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "DElete succes cok")
	})
	r.Patch("/testt", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Patch succes cok")
	})
	http.ListenAndServe(":8080", r)
}

// func Loggerscheck(f http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		w.Write([]byte("405 - Method Not Allowed"))
// 		// f.ServeHTTP(w, r)
// 	})
// }
