package router

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/toutaio/toutago/internal/di"
	"github.com/toutaio/toutago/pkg/touta"
)

func TestChiRouter_Routes(t *testing.T) {
	container := di.NewContainer()
	router := NewChiRouter(container)

	called := false
	handler := func(ctx touta.Context) error {
		called = true
		return ctx.String(200, "OK")
	}

	router.GET("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.Native().(*chi.Mux).ServeHTTP(w, req)

	if !called {
		t.Error("Handler should have been called")
	}

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestChiRouter_Methods(t *testing.T) {
	container := di.NewContainer()
	router := NewChiRouter(container)

	methods := map[string]func(string, touta.HandlerFunc){
		"GET":    router.GET,
		"POST":   router.POST,
		"PUT":    router.PUT,
		"DELETE": router.DELETE,
		"PATCH":  router.PATCH,
	}

	for method, registerFunc := range methods {
		t.Run(method, func(t *testing.T) {
			called := false
			registerFunc("/test", func(ctx touta.Context) error {
				called = true
				return nil
			})

			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()
			router.Native().(*chi.Mux).ServeHTTP(w, req)

			if !called {
				t.Errorf("%s handler should have been called", method)
			}
		})
	}
}

func TestContext_JSON(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	err := ctx.JSON(200, map[string]string{"status": "ok"})

	if err != nil {
		t.Fatalf("JSON failed: %v", err)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type should be application/json")
	}
}

func TestContext_String(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	err := ctx.String(200, "Hello World")

	if err != nil {
		t.Fatalf("String failed: %v", err)
	}

	if w.Body.String() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", w.Body.String())
	}

	if w.Header().Get("Content-Type") != "text/plain" {
		t.Error("Content-Type should be text/plain")
	}
}

func TestContext_HTML(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	err := ctx.HTML(200, "<h1>Hello</h1>")

	if err != nil {
		t.Fatalf("HTML failed: %v", err)
	}

	if w.Body.String() != "<h1>Hello</h1>" {
		t.Errorf("Expected '<h1>Hello</h1>', got '%s'", w.Body.String())
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("Content-Type should be text/html")
	}
}

func TestContext_GetSet(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	ctx.Set("key", "value")

	val := ctx.Get("key")
	if val != "value" {
		t.Errorf("Expected 'value', got '%v'", val)
	}
}

func TestContext_Query(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/?name=test", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	name := ctx.Query("name")

	if name != "test" {
		t.Errorf("Expected 'test', got '%s'", name)
	}
}
