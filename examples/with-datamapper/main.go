// Package main demonstrates using the DataMapper with Scéla message bus for data events.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
)

// User represents a user entity
type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserRepository manages user data operations
type UserRepository struct {
	bus touta.Bus
}

// Create creates a new user and publishes events
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	// Publish before-create event
	if err := r.bus.Publish(ctx, "user.creating", map[string]interface{}{
		"email": user.Email,
	}); err != nil {
		return err
	}
	
	// Simulate database insert
	user.ID = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	log.Printf("Created user: %+v", user)
	
	// Publish after-create event
	if err := r.bus.Publish(ctx, "user.created", map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	}); err != nil {
		return err
	}
	
	return nil
}

// Update updates a user and publishes events
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	// Publish before-update event
	if err := r.bus.Publish(ctx, "user.updating", map[string]interface{}{
		"id": user.ID,
	}); err != nil {
		return err
	}
	
	// Simulate database update
	user.UpdatedAt = time.Now()
	
	log.Printf("Updated user: %+v", user)
	
	// Publish after-update event
	if err := r.bus.Publish(ctx, "user.updated", map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	}); err != nil {
		return err
	}
	
	return nil
}

// Delete deletes a user and publishes events
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	// Publish before-delete event
	if err := r.bus.Publish(ctx, "user.deleting", map[string]interface{}{
		"id": id,
	}); err != nil {
		return err
	}
	
	// Simulate database delete
	log.Printf("Deleted user with ID: %d", id)
	
	// Publish after-delete event
	if err := r.bus.Publish(ctx, "user.deleted", map[string]interface{}{
		"id": id,
	}); err != nil {
		return err
	}
	
	return nil
}

func main() {
	// Create message bus
	bus := integration.NewScelaWithDefaults()
	defer bus.Close()
	
	// Register event handlers
	
	// Audit logger - logs all user events
	bus.Subscribe("user.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		log.Printf("[AUDIT] %s: %v", msg.Topic(), msg.Payload())
		return nil
	}))
	
	// Email notifications for user creation
	bus.Subscribe("user.created", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload := msg.Payload().(map[string]interface{})
		email := payload["email"].(string)
		name := payload["name"].(string)
		
		log.Printf("[EMAIL] Sending welcome email to %s (%s)", name, email)
		// In real app: send actual email
		return nil
	}))
	
	// Cache invalidation on updates
	bus.Subscribe("user.updated", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		payload := msg.Payload().(map[string]interface{})
		id := payload["id"].(int)
		
		log.Printf("[CACHE] Invalidating cache for user ID: %d", id)
		// In real app: invalidate cache
		return nil
	}))
	
	// Analytics tracking
	bus.Subscribe("user.created", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		log.Printf("[ANALYTICS] New user signup tracked")
		// In real app: send to analytics service
		return nil
	}))
	
	// Search index update
	bus.Subscribe("user.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		log.Printf("[SEARCH] Updating search index for event: %s", msg.Topic())
		// In real app: update search index
		return nil
	}))
	
	// Create repository
	repo := &UserRepository{bus: bus}
	
	// Demo: Create a user
	log.Println("\n=== Creating User ===")
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	
	ctx := context.Background()
	if err := repo.Create(ctx, user); err != nil {
		log.Fatal(err)
	}
	
	// Wait for async events to process
	time.Sleep(100 * time.Millisecond)
	
	// Demo: Update the user
	log.Println("\n=== Updating User ===")
	user.Name = "John Smith"
	if err := repo.Update(ctx, user); err != nil {
		log.Fatal(err)
	}
	
	// Wait for async events to process
	time.Sleep(100 * time.Millisecond)
	
	// Demo: Delete the user
	log.Println("\n=== Deleting User ===")
	if err := repo.Delete(ctx, user.ID); err != nil {
		log.Fatal(err)
	}
	
	// Wait for async events to process
	time.Sleep(100 * time.Millisecond)
	
	fmt.Println("\n=== Example Complete ===")
	fmt.Println("This example demonstrated:")
	fmt.Println("  - Event-driven data operations")
	fmt.Println("  - Multiple subscribers to same events")
	fmt.Println("  - Audit logging")
	fmt.Println("  - Email notifications")
	fmt.Println("  - Cache invalidation")
	fmt.Println("  - Analytics tracking")
	fmt.Println("  - Search index updates")
}
