// Package main demonstrates basic usage of the Toutā framework with integrated components.
package main

import (
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
	
	// Resolve logger
	loggerInterface, err := container.Make((*Logger)(nil))
	if err != nil {
		log.Fatal(err)
	}
	logger := loggerInterface.(Logger)
	
	logger.Info("Starting application with integrated components...")
	
	// Create router using cosan
	router := integration.NewRouter(container)
	
	// Register routes
	router.GET("/", func(ctx touta.Context) error {
		logger.Info("Handling GET /")
		return ctx.JSON(200, map[string]string{
			"message": "Welcome to Toutā with integrated components!",
			"version": "2.0.0",
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
