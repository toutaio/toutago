package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
)

// TestNewRouter verifies that NewRouter creates a valid router.
func TestNewRouter(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	if router == nil {
		t.Fatal("NewRouter returned nil")
	}
	
	// Verify it implements touta.Router interface
	var _ touta.Router = router
}

// TestRouterGET verifies GET route registration and handling.
func TestRouterGET(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	called := false
	router.GET("/test", func(ctx touta.Context) error {
		called = true
		return ctx.String(200, "OK")
	})
	
	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	// Serve the request using the native router
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if !called {
		t.Fatal("Handler was not called")
	}
	
	if w.Code != 200 {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

// TestRouterPOST verifies POST route registration.
func TestRouterPOST(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	called := false
	router.POST("/users", func(ctx touta.Context) error {
		called = true
		return ctx.JSON(201, map[string]string{"status": "created"})
	})
	
	req := httptest.NewRequest("POST", "/users", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if !called {
		t.Fatal("POST handler was not called")
	}
	
	if w.Code != 201 {
		t.Fatalf("Expected status 201, got %d", w.Code)
	}
}

// TestRouterGroup verifies route grouping functionality.
func TestRouterGroup(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	apiGroup := router.Group("/api")
	
	if apiGroup == nil {
		t.Fatal("Group returned nil")
	}
	
	// Verify group implements touta.Router
	var _ touta.Router = apiGroup
}

// TestRouterMiddleware verifies middleware registration.
func TestRouterMiddleware(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	middlewareCalled := false
	handlerCalled := false
	
	// Add middleware
	router.Use(func(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
		return func(ctx touta.Context) error {
			middlewareCalled = true
			return next(ctx)
		}
	})
	
	// Add handler
	router.GET("/test", func(ctx touta.Context) error {
		handlerCalled = true
		return ctx.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if !middlewareCalled {
		t.Fatal("Middleware was not called")
	}
	
	if !handlerCalled {
		t.Fatal("Handler was not called")
	}
}

// TestContextMethods verifies context adapter methods.
func TestContextMethods(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	var capturedCtx touta.Context
	
	router.GET("/context-test", func(ctx touta.Context) error {
		capturedCtx = ctx
		
		// Test Set/Get
		ctx.Set("key", "value")
		if val := ctx.Get("key"); val != "value" {
			t.Errorf("Get returned %v, expected 'value'", val)
		}
		
		// Test Container
		if ctx.Container() == nil {
			t.Error("Container() returned nil")
		}
		
		return nil
	})
	
	req := httptest.NewRequest("GET", "/context-test", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if capturedCtx == nil {
		t.Fatal("Context was not captured")
	}
}

// TestRouterPUT verifies PUT route registration.
func TestRouterPUT(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	called := false
	router.PUT("/users/1", func(ctx touta.Context) error {
		called = true
		return ctx.JSON(200, map[string]string{"status": "updated"})
	})
	
	req := httptest.NewRequest("PUT", "/users/1", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if !called {
		t.Fatal("PUT handler was not called")
	}
}

// TestRouterDELETE verifies DELETE route registration.
func TestRouterDELETE(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	called := false
	router.DELETE("/users/1", func(ctx touta.Context) error {
		called = true
		return ctx.String(204, "")
	})
	
	req := httptest.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if !called {
		t.Fatal("DELETE handler was not called")
	}
}

// TestRouterPATCH verifies PATCH route registration.
func TestRouterPATCH(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	called := false
	router.PATCH("/users/1", func(ctx touta.Context) error {
		called = true
		return ctx.JSON(200, map[string]string{"status": "patched"})
	})
	
	req := httptest.NewRequest("PATCH", "/users/1", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if !called {
		t.Fatal("PATCH handler was not called")
	}
}

// TestContextParam verifies parameter extraction.
func TestContextParam(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	var param string
	router.GET("/users/:id", func(ctx touta.Context) error {
		param = ctx.Param("id")
		return ctx.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if param != "123" {
		t.Errorf("Expected param '123', got '%s'", param)
	}
}

// TestContextQuery verifies query parameter extraction.
func TestContextQuery(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	var query string
	router.GET("/search", func(ctx touta.Context) error {
		query = ctx.Query("q")
		return ctx.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/search?q=test", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if query != "test" {
		t.Errorf("Expected query 'test', got '%s'", query)
	}
}

// TestContextRedirect verifies redirect functionality.
func TestContextRedirect(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	router.GET("/old", func(ctx touta.Context) error {
		return ctx.Redirect(302, "/new")
	})
	
	req := httptest.NewRequest("GET", "/old", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if w.Code != 302 {
		t.Errorf("Expected status 302, got %d", w.Code)
	}
	
	location := w.Header().Get("Location")
	if location != "/new" {
		t.Errorf("Expected location '/new', got '%s'", location)
	}
}

// TestContextRequestResponse verifies request/response access.
func TestContextRequestResponse(t *testing.T) {
	container := integration.NewContainer()
	router := integration.NewRouter(container)
	
	var hasRequest, hasResponse bool
	router.GET("/test", func(ctx touta.Context) error {
		hasRequest = ctx.Request() != nil
		hasResponse = ctx.Response() != nil
		return ctx.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	if nativeRouter, ok := router.Native().(http.Handler); ok {
		nativeRouter.ServeHTTP(w, req)
	}
	
	if !hasRequest {
		t.Error("Request() returned nil")
	}
	if !hasResponse {
		t.Error("Response() returned nil")
	}
}
