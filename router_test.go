package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {
	r := SetupDefaultRouter()
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Fatalf("unexpected body %q", w.Body.String())
	}
}

func TestGroupRouting(t *testing.T) {
	r := SetupDefaultRouter()
	g := r.Group("/api")
	g.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "pong" {
		t.Fatalf("unexpected response: code=%d body=%q", w.Code, w.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	r := SetupDefaultRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Test", "ok")
			next.ServeHTTP(w, req)
		})
	})

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hi")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Test") != "ok" {
		t.Fatalf("middleware not executed")
	}
}
