package integration

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
)

// TestRouterHooks tests BeforeRequest and AfterResponse hooks
func TestRouterHooks(t *testing.T) {
	container := NewContainer()
	router := NewRouter(container)

	var beforeCalled bool
	var afterCalled bool
	var afterStatusCode int

	// Register hooks
	router.BeforeRequest(func(req *http.Request) error {
		beforeCalled = true
		return nil
	})

	router.AfterResponse(func(req *http.Request, statusCode int) {
		afterCalled = true
		afterStatusCode = statusCode
	})

	router.GET("/test", func(ctx touta.Context) error {
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.Native().(http.Handler).ServeHTTP(w, req)

	if !beforeCalled {
		t.Error("BeforeRequest hook was not called")
	}

	if !afterCalled {
		t.Error("AfterResponse hook was not called")
	}

	if afterStatusCode != 200 {
		t.Errorf("Expected status code 200 in AfterResponse, got %d", afterStatusCode)
	}
}

// TestRouterHooksAbort tests that BeforeRequest can abort requests
func TestRouterHooksAbort(t *testing.T) {
	container := NewContainer()
	router := NewRouter(container)

	var handlerCalled bool

	router.BeforeRequest(func(req *http.Request) error {
		return errors.New("request aborted")
	})

	router.GET("/test", func(ctx touta.Context) error {
		handlerCalled = true
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.Native().(http.Handler).ServeHTTP(w, req)

	if handlerCalled {
		t.Error("Handler should not have been called after hook abort")
	}

	if w.Code != 500 {
		t.Errorf("Expected status 500 after hook abort, got %d", w.Code)
	}
}

// TestRouterErrorHandler tests custom error handling
func TestRouterErrorHandler(t *testing.T) {
	container := NewContainer()
	router := NewRouter(container)

	var errorHandlerCalled bool
	var capturedError error

	router.SetErrorHandler(func(ctx touta.Context, err error) {
		errorHandlerCalled = true
		capturedError = err
		ctx.JSON(418, map[string]string{"error": err.Error()})
	})

	router.GET("/error", func(ctx touta.Context) error {
		return errors.New("test error")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	router.Native().(http.Handler).ServeHTTP(w, req)

	if !errorHandlerCalled {
		t.Error("Error handler was not called")
	}

	if capturedError == nil || capturedError.Error() != "test error" {
		t.Errorf("Expected error 'test error', got %v", capturedError)
	}

	if w.Code != 418 {
		t.Errorf("Expected custom status 418, got %d", w.Code)
	}
}

// TestRouterIntrospection tests GetRoutes and FindRoute
func TestRouterIntrospection(t *testing.T) {
	container := NewContainer()
	router := NewRouter(container)

	router.GET("/users", func(ctx touta.Context) error {
		return ctx.JSON(200, nil)
	})

	router.POST("/users", func(ctx touta.Context) error {
		return ctx.JSON(201, nil)
	})

	router.GET("/users/:id", func(ctx touta.Context) error {
		return ctx.JSON(200, nil)
	})

	routes := router.GetRoutes()

	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Check that routes contain expected methods
	methodCount := make(map[string]int)
	for _, route := range routes {
		methodCount[route.Method]++
	}

	if methodCount["GET"] != 2 {
		t.Errorf("Expected 2 GET routes, got %d", methodCount["GET"])
	}

	if methodCount["POST"] != 1 {
		t.Errorf("Expected 1 POST route, got %d", methodCount["POST"])
	}
}

// TestRouterHooksOrder tests that multiple hooks execute in order
func TestRouterHooksOrder(t *testing.T) {
	container := NewContainer()
	router := NewRouter(container)

	var executionOrder []string

	router.BeforeRequest(func(req *http.Request) error {
		executionOrder = append(executionOrder, "before1")
		return nil
	})

	router.BeforeRequest(func(req *http.Request) error {
		executionOrder = append(executionOrder, "before2")
		return nil
	})

	router.AfterResponse(func(req *http.Request, statusCode int) {
		executionOrder = append(executionOrder, "after1")
	})

	router.AfterResponse(func(req *http.Request, statusCode int) {
		executionOrder = append(executionOrder, "after2")
	})

	router.GET("/test", func(ctx touta.Context) error {
		executionOrder = append(executionOrder, "handler")
		return ctx.JSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.Native().(http.Handler).ServeHTTP(w, req)

	expected := []string{"before1", "before2", "handler", "after1", "after2"}
	if len(executionOrder) != len(expected) {
		t.Errorf("Expected %d executions, got %d", len(expected), len(executionOrder))
	}

	for i, step := range expected {
		if i >= len(executionOrder) || executionOrder[i] != step {
			t.Errorf("Expected step %d to be '%s', got '%s'", i, step, executionOrder[i])
		}
	}
}
