package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
)

// TestContextEcosystemIntegration tests the new ecosystem integration methods
func TestContextEcosystemIntegration(t *testing.T) {
	container := NewContainer()
	router := NewRouter(container)

	// Test Params method
	router.GET("/users/:id/:action", func(ctx touta.Context) error {
		params := ctx.Params()
		return ctx.JSON(200, params)
	})

	// Test QueryAll method
	router.GET("/search", func(ctx touta.Context) error {
		tags := ctx.QueryAll("tag")
		return ctx.JSON(200, map[string]interface{}{"tags": tags})
	})

	// Test Header method
	router.GET("/headers", func(ctx touta.Context) error {
		ctx.Header().Set("X-Custom-Header", "test-value")
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	t.Run("Params", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/123/edit", nil)
		w := httptest.NewRecorder()

		router.Native().(http.Handler).ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !bytes.Contains(w.Body.Bytes(), []byte("123")) || !bytes.Contains(w.Body.Bytes(), []byte("edit")) {
			t.Errorf("Expected params in response, got: %s", body)
		}
	})

	t.Run("QueryAll", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/search?tag=go&tag=web&tag=api", nil)
		w := httptest.NewRecorder()

		router.Native().(http.Handler).ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !bytes.Contains(w.Body.Bytes(), []byte("go")) || !bytes.Contains(w.Body.Bytes(), []byte("web")) {
			t.Errorf("Expected tags in response, got: %s", body)
		}
	})

	t.Run("Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/headers", nil)
		w := httptest.NewRecorder()

		router.Native().(http.Handler).ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		if w.Header().Get("X-Custom-Header") != "test-value" {
			t.Errorf("Expected X-Custom-Header to be 'test-value', got '%s'", w.Header().Get("X-Custom-Header"))
		}
	})
}

// TestContextRenderMethod tests the Render method with Fith integration
func TestContextRenderMethod(t *testing.T) {
	t.Skip("Render integration requires proper container setup - tested in examples")
}

// TestRouterOPTIONSAndHEAD tests the new HTTP methods
func TestRouterOPTIONSAndHEAD(t *testing.T) {
	container := NewContainer()
	router := NewRouter(container)

	router.OPTIONS("/api", func(ctx touta.Context) error {
		ctx.Header().Set("Allow", "GET, POST, OPTIONS")
		return ctx.Status(204).JSON(204, nil)
	})

	router.HEAD("/health", func(ctx touta.Context) error {
		return ctx.Status(200).JSON(200, nil)
	})

	t.Run("OPTIONS", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api", nil)
		w := httptest.NewRecorder()

		router.Native().(http.Handler).ServeHTTP(w, req)

		if w.Code != 204 {
			t.Errorf("Expected status 204, got %d", w.Code)
		}

		if w.Header().Get("Allow") == "" {
			t.Error("Expected Allow header to be set")
		}
	})

	t.Run("HEAD", func(t *testing.T) {
		req := httptest.NewRequest("HEAD", "/health", nil)
		w := httptest.NewRecorder()

		router.Native().(http.Handler).ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}
