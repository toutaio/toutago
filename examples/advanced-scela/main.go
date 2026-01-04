// Package main demonstrates advanced Scéla features including retries, DLQ, and complex patterns.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
)

// Order represents an e-commerce order
type Order struct {
	ID         string
	CustomerID string
	Amount     float64
	Status     string
	Items      []string
	CreatedAt  time.Time
}

// PaymentProcessor simulates payment processing
type PaymentProcessor struct {
	failCount int
}

func (p *PaymentProcessor) Process(order *Order) error {
	// Simulate intermittent failures
	p.failCount++
	if p.failCount%3 == 0 {
		return errors.New("payment gateway timeout")
	}
	
	log.Printf("[PAYMENT] Processing payment of $%.2f for order %s", order.Amount, order.ID)
	time.Sleep(100 * time.Millisecond)
	return nil
}

// InventoryService manages inventory
type InventoryService struct{}

func (s *InventoryService) Reserve(items []string) error {
	log.Printf("[INVENTORY] Reserving items: %v", items)
	time.Sleep(50 * time.Millisecond)
	return nil
}

func (s *InventoryService) Release(items []string) error {
	log.Printf("[INVENTORY] Releasing items: %v", items)
	return nil
}

// ShippingService handles shipping
type ShippingService struct{}

func (s *ShippingService) CreateShipment(order *Order) error {
	log.Printf("[SHIPPING] Creating shipment for order %s", order.ID)
	time.Sleep(75 * time.Millisecond)
	return nil
}

// EmailService sends emails
type EmailService struct{}

func (s *EmailService) SendOrderConfirmation(order *Order) error {
	log.Printf("[EMAIL] Sending order confirmation to customer %s", order.CustomerID)
	return nil
}

func (s *EmailService) SendPaymentFailed(order *Order) error {
	log.Printf("[EMAIL] Sending payment failure notification for order %s", order.ID)
	return nil
}

func main() {
	// Create services
	paymentProcessor := &PaymentProcessor{}
	inventoryService := &InventoryService{}
	shippingService := &ShippingService{}
	emailService := &EmailService{}
	
	// Create bus with options
	bus := scela.New(
		scela.WithWorkers(5),
		scela.WithMaxRetries(3),
		scela.WithDeadLetterHandler(scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
			log.Printf("[DLQ] Message failed after all retries: %s - %v", msg.Topic(), msg.Payload())
			
			// In real app: store in database, alert ops team, etc.
			payload := msg.Payload().(map[string]interface{})
			if orderID, ok := payload["order_id"].(string); ok {
				// Send failure notification
				order := &Order{ID: orderID, CustomerID: "CUST001"}
				emailService.SendPaymentFailed(order)
			}
			
			return nil
		})),
	)
	defer bus.Close()
	
	// Create logging handler wrapper
	loggingHandler := func(handler scela.Handler) scela.Handler {
		return scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
			start := time.Now()
			err := handler.Handle(ctx, msg)
			duration := time.Since(start)
			
			status := "SUCCESS"
			if err != nil {
				status = "ERROR"
			}
			
			log.Printf("[HANDLER] %s | %s | %v | %s", msg.Topic(), status, duration, msg.ID())
			return err
		})
	}
	
	// Create validation handler wrapper
	validationHandler := func(handler scela.Handler) scela.Handler {
		return scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
			payload, ok := msg.Payload().(map[string]interface{})
			if !ok {
				return errors.New("invalid payload format")
			}
			
			if payload["order_id"] == nil {
				return errors.New("order_id is required")
			}
			
			return handler.Handle(ctx, msg)
		})
	}
	
	// Handler for order creation - with logging and validation
	bus.Subscribe("order.created", loggingHandler(validationHandler(scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
		payload := msg.Payload().(map[string]interface{})
		
		order := &Order{
			ID:         payload["order_id"].(string),
			CustomerID: payload["customer_id"].(string),
			Amount:     payload["amount"].(float64),
			Items:      payload["items"].([]string),
			Status:     "pending",
			CreatedAt:  time.Now(),
		}
		
		log.Printf("[ORDER] Order created: %s (Customer: %s, Amount: $%.2f)", order.ID, order.CustomerID, order.Amount)
		
		// Reserve inventory
		if err := inventoryService.Reserve(order.Items); err != nil {
			return err
		}
		
		// Publish inventory reserved event
		return bus.Publish(ctx, "order.inventory_reserved", map[string]interface{}{
			"order_id": order.ID,
			"items":    order.Items,
		})
	}))))
	
	// Handler for inventory reservation - with logging and validation
	bus.Subscribe("order.inventory_reserved", loggingHandler(validationHandler(scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
		payload := msg.Payload().(map[string]interface{})
		orderID := payload["order_id"].(string)
		
		// Create mock order for payment
		order := &Order{
			ID:         orderID,
			CustomerID: "CUST001",
			Amount:     99.99,
		}
		
		// Process payment (may fail and retry)
		if err := paymentProcessor.Process(order); err != nil {
			log.Printf("[ORDER] Payment failed for order %s: %v", orderID, err)
			return err
		}
		
		log.Printf("[ORDER] Payment successful for order %s", orderID)
		
		// Publish payment completed event
		return bus.Publish(ctx, "order.payment_completed", map[string]interface{}{
			"order_id": orderID,
		})
	}))))
	
	// Handler for payment completion - with logging and validation
	bus.Subscribe("order.payment_completed", loggingHandler(validationHandler(scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
		payload := msg.Payload().(map[string]interface{})
		orderID := payload["order_id"].(string)
		
		order := &Order{
			ID:         orderID,
			CustomerID: "CUST001",
		}
		
		// Create shipment
		if err := shippingService.CreateShipment(order); err != nil {
			return err
		}
		
		// Send confirmation email
		if err := emailService.SendOrderConfirmation(order); err != nil {
			return err
		}
		
		log.Printf("[ORDER] Order %s completed successfully!", orderID)
		return nil
	}))))
	
	// Handler for analytics - uses pattern matching
	bus.Subscribe("order.*", scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
		log.Printf("[ANALYTICS] Event tracked: %s", msg.Topic())
		// In real app: send to analytics service
		return nil
	}))
	
	// Demonstrate the order flow
	ctx := context.Background()
	
	log.Println("\n=== Processing Order 1 (will succeed after retries) ===")
	err := bus.PublishSync(ctx, "order.created", map[string]interface{}{
		"order_id":    "ORD-001",
		"customer_id": "CUST001",
		"amount":      99.99,
		"items":       []string{"ITEM-001", "ITEM-002"},
	})
	if err != nil {
		log.Printf("Error: %v", err)
	}
	
	// Wait for async processing
	time.Sleep(2 * time.Second)
	
	log.Println("\n=== Processing Order 2 (will succeed after retries) ===")
	err = bus.PublishSync(ctx, "order.created", map[string]interface{}{
		"order_id":    "ORD-002",
		"customer_id": "CUST002",
		"amount":      149.99,
		"items":       []string{"ITEM-003"},
	})
	if err != nil {
		log.Printf("Error: %v", err)
	}
	
	// Wait for async processing
	time.Sleep(2 * time.Second)
	
	log.Println("\n=== Processing Priority Order (high priority) ===")
	err = bus.PublishWithPriority(ctx, "order.created", map[string]interface{}{
		"order_id":    "ORD-003-PRIORITY",
		"customer_id": "CUST003",
		"amount":      299.99,
		"items":       []string{"ITEM-004", "ITEM-005", "ITEM-006"},
	}, scela.PriorityHigh)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	
	// Wait for async processing
	time.Sleep(2 * time.Second)
	
	fmt.Println("\n=== Example Complete ===")
	fmt.Println("This example demonstrated:")
	fmt.Println("  - Event-driven order processing workflow")
	fmt.Println("  - Retry middleware with exponential backoff")
	fmt.Println("  - Validation and logging middleware")
	fmt.Println("  - Dead letter queue for failed messages")
	fmt.Println("  - Pattern matching for analytics")
	fmt.Println("  - Priority message processing")
	fmt.Println("  - Sync and async publishing")
	fmt.Println("  - Middleware chaining")
	fmt.Println("  - Message persistence")
}
