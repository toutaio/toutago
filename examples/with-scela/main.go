package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
	"github.com/toutaio/toutago-scela-bus/pkg/scela"
)

// LoggingMiddleware logs all messages
func LoggingMiddleware(next touta.Handler) touta.Handler {
	return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		start := time.Now()
		fmt.Printf("[MIDDLEWARE] Processing %s...\n", msg.Topic())
		err := next.Handle(ctx, msg)
		fmt.Printf("[MIDDLEWARE] Completed %s in %v (error: %v)\n", msg.Topic(), time.Since(start), err)
		return err
	})
}

// ValidationMiddleware validates message payload
func ValidationMiddleware(next touta.Handler) touta.Handler {
	return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload, ok := msg.Payload().(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid payload type")
		}
		if _, ok := payload["id"]; !ok {
			return fmt.Errorf("missing required field: id")
		}
		return next.Handle(ctx, msg)
	})
}

func main() {
	// Create bus with middleware
	bus := integration.NewScelaBusWithMiddleware(LoggingMiddleware)
	defer bus.Close()

	// Subscribe to all user events with pattern matching
	_, err := bus.Subscribe("user.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		fmt.Printf("[HANDLER] User event: %s\n", msg.Topic())
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to user.created with validation middleware
	handler := ValidationMiddleware(touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload := msg.Payload().(map[string]interface{})
		fmt.Printf("[HANDLER] New user created: %v\n", payload)
		
		// Simulate sending welcome email
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[HANDLER] Welcome email sent to user %s\n", payload["id"])
		return nil
	}))
	
	_, err = bus.Subscribe("user.created", handler)
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to user.deleted
	_, err = bus.Subscribe("user.deleted", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload := msg.Payload().(map[string]interface{})
		fmt.Printf("[HANDLER] User deleted: %v\n", payload)
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to high-priority notifications
	_, err = bus.Subscribe("notification.urgent", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		fmt.Printf("[URGENT] %v\n", msg.Payload())
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}

	// HTTP handlers
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userData := map[string]interface{}{
				"id":    "123",
				"email": "user@example.com",
				"name":  "John Doe",
			}
			
			// Publish async (non-blocking)
			if err := bus.Publish(r.Context(), "user.created", userData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "User created\n")

		case http.MethodDelete:
			if err := bus.Publish(r.Context(), "user.deleted", map[string]interface{}{
				"id": "123",
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "User deleted\n")

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Publish high-priority notification
		if err := bus.PublishWithPriority(r.Context(), "notification.urgent", map[string]interface{}{
			"message": "Critical system alert!",
			"time":    time.Now(),
		}, scela.PriorityHigh); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Notification sent\n")
	})

	http.HandleFunc("/sync-event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Publish sync (wait for handlers to complete)
		if err := bus.PublishSync(r.Context(), "user.updated", map[string]interface{}{
			"id":   "123",
			"name": "Jane Doe",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Synchronous event completed\n")
	})

	fmt.Println("=== Scéla Message Bus Demo ===")
	fmt.Println("Server starting on :8080")
	fmt.Println()
	fmt.Println("Try these commands:")
	fmt.Println("  curl -X POST http://localhost:8080/users")
	fmt.Println("  curl -X DELETE http://localhost:8080/users")
	fmt.Println("  curl -X POST http://localhost:8080/notify")
	fmt.Println("  curl -X POST http://localhost:8080/sync-event")
	fmt.Println()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
