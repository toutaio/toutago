// Package main demonstrates using the Fíth template renderer with Scéla message bus for event-driven rendering.
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
)

// PageView tracks page view analytics
type PageView struct {
	Path      string
	Timestamp time.Time
	UserAgent string
}

// ViewTracker tracks all page views
type ViewTracker struct {
	mu    sync.RWMutex
	views []PageView
}

func (vt *ViewTracker) Add(view PageView) {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vt.views = append(vt.views, view)
}

func (vt *ViewTracker) GetStats() map[string]int {
	vt.mu.RLock()
	defer vt.mu.RUnlock()
	
	stats := make(map[string]int)
	for _, view := range vt.views {
		stats[view.Path]++
	}
	return stats
}

// SimpleRenderer is a basic template renderer for demonstration
type SimpleRenderer struct{}

func (r *SimpleRenderer) Render(template string, data interface{}) ([]byte, error) {
	dataMap := data.(map[string]interface{})
	
	var html string
	switch template {
	case "home":
		html = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>%s</title></head>
<body>
<h1>%s</h1>
<h2>Hello, %s!</h2>
<ul>`, dataMap["Title"], dataMap["Title"], dataMap["Name"])
		
		for _, feature := range dataMap["Features"].([]string) {
			html += fmt.Sprintf("<li>%s</li>", feature)
		}
		html += "</ul></body></html>"
		
	case "stats":
		html = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>View Statistics</title></head>
<body>
<h1>Page View Statistics</h1>
<ul>`)
		
		stats := dataMap["Stats"].(map[string]int)
		for path, count := range stats {
			html += fmt.Sprintf("<li>%s: %d views</li>", path, count)
		}
		html += "</ul></body></html>"
	}
	
	return []byte(html), nil
}

func main() {
	// Create DI container
	container := integration.NewContainer()
	
	// Create message bus
	bus := integration.NewScelaWithDefaults()
	defer bus.Close()
	
	// Create view tracker
	tracker := &ViewTracker{views: []PageView{}}
	
	// Subscribe to page view events for analytics
	bus.Subscribe("page.view", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload := msg.Payload().(map[string]interface{})
		view := PageView{
			Path:      payload["path"].(string),
			Timestamp: time.Now(),
			UserAgent: payload["user_agent"].(string),
		}
		tracker.Add(view)
		log.Printf("[ANALYTICS] Page view: %s", view.Path)
		return nil
	}))
	
	// Subscribe to template rendering events for caching
	bus.Subscribe("template.render", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload := msg.Payload().(map[string]interface{})
		template := payload["template"].(string)
		log.Printf("[CACHE] Template rendered: %s", template)
		// In real app: cache the rendered output
		return nil
	}))
	
	// Subscribe to slow render events for monitoring
	bus.Subscribe("template.render.slow", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload := msg.Payload().(map[string]interface{})
		template := payload["template"].(string)
		duration := payload["duration"].(time.Duration)
		log.Printf("[MONITOR] Slow render detected: %s took %v", template, duration)
		return nil
	}))
	
	// Create simple renderer
	renderer := &SimpleRenderer{}
	
	// Bind renderer to container
	container.Singleton((*touta.TemplateRenderer)(nil), renderer)
	
	// Create router
	router := integration.NewRouter(container)
	
	// Middleware to publish page view events
	router.Use(func(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
		return func(ctx touta.Context) error {
			// Publish page view event
			bus.Publish(context.Background(), "page.view", map[string]interface{}{
				"path":       ctx.Request().URL.Path,
				"user_agent": ctx.Request().UserAgent(),
			})
			return next(ctx)
		}
	})
	
	// Home page with template
	router.GET("/", func(ctx touta.Context) error {
		start := time.Now()
		
		data := map[string]interface{}{
			"Title": "Welcome to Toutā",
			"Name":  "Developer",
			"Features": []string{
				"Integrated component architecture",
				"Fíth template engine",
				"Cosan HTTP router",
				"Nasc dependency injection",
				"Scéla message bus",
			},
		}
		
		output, err := renderer.Render("home", data)
		if err != nil {
			return err
		}
		
		duration := time.Since(start)
		
		// Publish render event
		bus.Publish(context.Background(), "template.render", map[string]interface{}{
			"template": "home",
			"duration": duration,
		})
		
		// If render was slow, publish slow render event
		if duration > 100*time.Millisecond {
			bus.Publish(context.Background(), "template.render.slow", map[string]interface{}{
				"template": "home",
				"duration": duration,
			})
		}
		
		return ctx.HTML(200, string(output))
	})
	
	// Stats page showing analytics
	router.GET("/stats", func(ctx touta.Context) error {
		stats := tracker.GetStats()
		
		data := map[string]interface{}{
			"Stats": stats,
		}
		
		output, err := renderer.Render("stats", data)
		if err != nil {
			return err
		}
		
		return ctx.HTML(200, string(output))
	})
	
	log.Println("Server starting on :8080")
	log.Println("Visit http://localhost:8080/ for the home page")
	log.Println("Visit http://localhost:8080/stats for analytics")
	log.Fatal(router.Listen(":8080"))
}
