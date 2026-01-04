package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
)

// BenchmarkRouterSimpleRoute benchmarks a simple route lookup
func BenchmarkRouterSimpleRoute(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	router.GET("/users", func(ctx touta.Context) error {
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkRouterParametricRoute benchmarks a parametric route
func BenchmarkRouterParametricRoute(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	router.GET("/users/:id", func(ctx touta.Context) error {
		id := ctx.Param("id")
		return ctx.JSON(200, map[string]string{"id": id})
	})

	req := httptest.NewRequest("GET", "/users/123", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkRouterMultiParam benchmarks routes with multiple parameters
func BenchmarkRouterMultiParam(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	router.GET("/users/:id/posts/:postId", func(ctx touta.Context) error {
		return ctx.JSON(200, ctx.Params())
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkRouterManyRoutes benchmarks routing with many registered routes
func BenchmarkRouterManyRoutes(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	// Register 100 routes
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/route%d", i)
		router.GET(path, func(ctx touta.Context) error {
			return ctx.JSON(200, nil)
		})
	}

	req := httptest.NewRequest("GET", "/route50", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkRouterWithMiddleware benchmarks routing with middleware
func BenchmarkRouterWithMiddleware(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	// Add middleware
	router.Use(func(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
		return func(ctx touta.Context) error {
			return next(ctx)
		}
	})

	router.GET("/users", func(ctx touta.Context) error {
		return ctx.JSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkRouterWithHooks benchmarks routing with before/after hooks
func BenchmarkRouterWithHooks(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	var counter int
	router.BeforeRequest(func(req *http.Request) error {
		counter++
		return nil
	})

	router.AfterResponse(func(req *http.Request, statusCode int) {
		counter++
	})

	router.GET("/users", func(ctx touta.Context) error {
		return ctx.JSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkContextParamAccess benchmarks parameter access
func BenchmarkContextParamAccess(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	router.GET("/users/:id", func(ctx touta.Context) error {
		_ = ctx.Param("id")
		return nil
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkContextJSONResponse benchmarks JSON response writing
func BenchmarkContextJSONResponse(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	data := map[string]interface{}{
		"id":      123,
		"name":    "Test User",
		"email":   "test@example.com",
		"active":  true,
		"credits": 100,
	}

	router.GET("/users", func(ctx touta.Context) error {
		return ctx.JSON(200, data)
	})

	req := httptest.NewRequest("GET", "/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}

// BenchmarkRouterGrouping benchmarks route groups
func BenchmarkRouterGrouping(b *testing.B) {
	container := NewContainer()
	router := NewRouter(container)

	api := router.Group("/api/v1")
	api.GET("/users", func(ctx touta.Context) error {
		return ctx.JSON(200, nil)
	})

	req := httptest.NewRequest("GET", "/api/v1/users", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.Native().(http.Handler).ServeHTTP(w, req)
	}
}
