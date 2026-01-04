// Package main demonstrates basic usage of the Toutā framework with integrated components.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
)

// Logger is a simple logging interface
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// ConsoleLogger logs to stdout
type ConsoleLogger struct {
	prefix string
}

func (l *ConsoleLogger) Info(msg string) {
	fmt.Printf("[INFO] %s: %s\n", l.prefix, msg)
}

func (l *ConsoleLogger) Error(msg string) {
	fmt.Printf("[ERROR] %s: %s\n", l.prefix, msg)
}

func main() {
	// Create DI container using nasc
	container := integration.NewContainer()
	
	// Bind services
	err := container.Singleton((*Logger)(nil), &ConsoleLogger{prefix: "app"})
	if err != nil {
		log.Fatal(err)
	}
	
	// Create and register message bus using Scéla
	bus := integration.NewScelaWithDefaults()
	defer bus.Close()
	
	err = container.Factory((*touta.Bus)(nil), func(c touta.Container) (interface{}, error) {
		return bus, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Resolve logger
	loggerInterface, err := container.Make((*Logger)(nil))
	if err != nil {
		log.Fatal(err)
	}
	logger := loggerInterface.(Logger)
	
	logger.Info("Starting application with integrated components...")
	
	// Subscribe to application events
	bus.Subscribe("app.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		logger.Info(fmt.Sprintf("Event: %s - %v", msg.Topic(), msg.Payload()))
		return nil
	}))
	
	// Publish application startup event
	ctx := context.Background()
	bus.Publish(ctx, "app.started", map[string]interface{}{
		"version": "2.0.0",
		"components": []string{"nasc", "cosan", "scela"},
	})
	
	// Create router using cosan
	router := integration.NewRouter(container)
	
	// Register routes
	router.GET("/", func(ctx touta.Context) error {
		logger.Info("Handling GET /")
		
		// Publish request event
		bus.Publish(context.Background(), "app.request", map[string]interface{}{
			"path": "/",
			"method": "GET",
		})
		
		return ctx.JSON(200, map[string]string{
			"message": "Welcome to Toutā with integrated components!",
			"version": "2.0.0",
			"components": "nasc + cosan + scela",
		})
	})
	
	router.GET("/health", func(ctx touta.Context) error {
		return ctx.JSON(200, map[string]string{
			"status": "healthy",
		})
	})
	
	// Add middleware
	router.Use(func(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
		return func(ctx touta.Context) error {
			logger.Info(fmt.Sprintf("%s %s", ctx.Request().Method, ctx.Request().URL.Path))
			return next(ctx)
		}
	})
	
	logger.Info("Server starting on :8080")
	log.Fatal(router.Listen(":8080"))
}
