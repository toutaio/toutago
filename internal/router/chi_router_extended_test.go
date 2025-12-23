package router

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/toutaio/toutago/internal/di"
	"github.com/toutaio/toutago/pkg/touta"
)

func TestChiRouter_Group(t *testing.T) {
	container := di.NewContainer()
	router := NewChiRouter(container)

	apiGroup := router.Group("/api")
	apiGroup.GET("/users", func(ctx touta.Context) error {
		return ctx.String(200, "users")
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()
	router.Native().(*chi.Mux).ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestChiRouter_Middleware(t *testing.T) {
	container := di.NewContainer()
	router := NewChiRouter(container)

	middlewareCalled := false
	middleware := func(next touta.HandlerFunc) touta.HandlerFunc {
		return func(ctx touta.Context) error {
			middlewareCalled = true
			return next(ctx)
		}
	}

	router.Use(middleware)
	router.GET("/test", func(ctx touta.Context) error {
		return ctx.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.Native().(*chi.Mux).ServeHTTP(w, req)

	if !middlewareCalled {
		t.Error("Middleware should have been called")
	}
}

func TestChiRouter_HandlerError(t *testing.T) {
	container := di.NewContainer()
	router := NewChiRouter(container)

	router.GET("/error", func(ctx touta.Context) error {
		return errors.New("test error")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.Native().(*chi.Mux).ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected 500 for error, got %d", w.Code)
	}
}

func TestContext_Redirect(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	err := ctx.Redirect(302, "/login")

	if err != nil {
		t.Fatalf("Redirect failed: %v", err)
	}

	if w.Code != 302 {
		t.Errorf("Expected 302, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if location != "/login" {
		t.Errorf("Expected Location /login, got %s", location)
	}
}

func TestContext_Status(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	result := ctx.Status(201)

	if result != ctx {
		t.Error("Status should return context for chaining")
	}

	if w.Code != 201 {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestContext_Request(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)

	if ctx.Request() != req {
		t.Error("Request() should return original request")
	}
}

func TestContext_Response(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)

	if ctx.Response() != w {
		t.Error("Response() should return response writer")
	}
}

func TestContext_Container(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)

	if ctx.Container() != container {
		t.Error("Container() should return DI container")
	}
}

func TestContext_ParamNotFound(t *testing.T) {
	container := di.NewContainer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req, container)
	param := ctx.Param("nonexistent")

	if param != "" {
		t.Error("Nonexistent param should return empty string")
	}
}

func TestChiRouter_Native(t *testing.T) {
	container := di.NewContainer()
	router := NewChiRouter(container)

	native := router.Native()
	if native == nil {
		t.Error("Native() should return underlying router")
	}

	_, ok := native.(*chi.Mux)
	if !ok {
		t.Error("Native router should be Chi Mux")
	}
}
