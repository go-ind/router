// func logRequest(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("[%s] %s\n", r.Method, r.URL.Path)
// 		next(w, r)
// 	}
// }

// func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Authentication logic here
// 		next(w, r)
// 	}
// }

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type router struct {
	routes      map[string]map[string]func(http.ResponseWriter, *http.Request)
	SettingLogs string
}

func (r *router) addRoute(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, ok := r.routes[path]; !ok {
		r.routes[path] = make(map[string]func(http.ResponseWriter, *http.Request))
	}
	r.routes[path][method] = handler
}
func (r *router) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodGet, path, handler)
}
func (r *router) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPost, path, handler)
}

func (r *router) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPut, path, handler)
}

func (r *router) Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodPatch, path, handler)
}

func (r *router) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.addRoute(HttpMethodDelete, path, handler)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	start := time.Now().In(loc)

	req = StartRecord(req, start)

	if handlers, ok := r.routes[req.URL.Path]; ok {
		if handler, ok := handlers[req.Method]; ok {
			handler(w, req)
			return
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
	http.NotFound(w, req)
}

func StartRecord(req *http.Request, start time.Time) *http.Request {
	ctx := req.Context()

	v := new(Data)
	v.RequestID = uuid.New().String()

	v.Host = req.Host
	v.Endpoint = req.URL.Path
	v.TimeStart = start
	v.Device = "Web-Base"

	v.RequestMethod = req.Method
	v.RequestHeader = DumpRequest(req)

	ctx = context.WithValue(ctx, LogKey, v)

	return req.WithContext(ctx)
}

func DumpRequest(req *http.Request) string {
	header, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "cannot dump request"
	}

	trim := bytes.ReplaceAll(header, []byte("\r\n"), []byte("   "))
	return string(trim)
}

func Response(w http.ResponseWriter, ctx context.Context, code int, status bool, message string, rs, pagination interface{}) {
	resservice := Responseservice{}
	resservice.Status = code
	if status {
		resservice.Data = rs
		resservice.Pagination = pagination
	} else {
		resservice.ErrorMessage = "Error"
	}

	var input []byte

	resservice.Message = message
	switch rs.(type) {
	case string:
		input = []byte(rs.(string))
	case []byte:
		input = rs.([]byte)
	default:
		input, _ = JSONMarshal(rs)
	}
	if ctx == nil {
		// Handle For CTX if is null
		ctx = context.TODO()
	}
	Logger(ctx, string(input), code)
	origin := "*"

	v, ok := ctx.Value(LogKey).(*Data)
	if ok {
		words := strings.Fields(v.RequestHeader)
		for i := 0; i < len(words); i++ {
			if words[i] == "Origin:" {
				origin = words[i+1]
				break
			}
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Strict-Transport-Security", "max-age=15552000; includeSubDomains")
	w.Header().Set("X-DNS-Prefetch-Control", "off")
	w.Header().Set("Vary", "X-HTTP-Method-Override")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resservice)
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
	logrus.SetFormatter(UTCFormatter{&logrus.JSONFormatter{}})
	if level == "ERROR" {
		logrus.WithField("data", out).Error("apps")
	} else if level == "INFO" {
		logrus.WithField("data", out).Info("apps")
	} else if level == "WARN" {
		logrus.WithField("data", out).Warn("apps")
	}
}

// UTCFormatter ...
type UTCFormatter struct {
	logrus.Formatter
}

// func (r *router) Group(prefix string, middleware ...func(http.HandlerFunc) http.HandlerFunc) *group {
// 	// newPrefix := g.prefix + prefix
// 	// handlerChain := chain(g.handlers, nil)
// 	// return &group{
// 	// 	router:   g.router,
// 	// 	prefix:   newPrefix,
// 	// 	handlers: append(middleware, handlerChain),
// 	// }
// 	// return &group{}
// }

func main() {
	r := &router{
		routes:      make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
		SettingLogs: "FULLY_LOGING",
	}

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Hello, world!")
	})
	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		// id := req.URL.Query().Get("id")
		// fmt.Println("id =>", id)
		// fmt.Fprint(w, "POST succes cok")
		Response(w, ctx, 200, true, "Succes", nil, nil)
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
	r.Patch("/", func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		fmt.Println("id =>", id)
		fmt.Fprint(w, "Patch succes cok")
	})

	http.ListenAndServe(":8080", r)
}
