package app

import (
	"context"
	"log"
	"sync"
)

// ShutdownCoordinator defines the contract for coordinating graceful shutdown
// of background services. This follows the Dependency Inversion Principle by
// depending on an abstraction rather than concrete implementations.
type ShutdownCoordinator interface {
	// Start launches a background service with proper shutdown coordination.
	// This follows the Single Responsibility Principle by separating lifecycle management.
	Start(ctx context.Context, name string, fn func(ctx context.Context) error)
	// Shutdown waits for all registered services to complete their shutdown.
	Shutdown(ctx context.Context) error
}

// WaitGroupCoordinator implements ShutdownCoordinator using sync.WaitGroup.
// This provides a concrete implementation while allowing for other implementations
// (e.g., for testing or different coordination strategies).
type WaitGroupCoordinator struct {
	wg sync.WaitGroup
}

// NewWaitGroupCoordinator creates a new WaitGroup-based shutdown coordinator.
func NewWaitGroupCoordinator() *WaitGroupCoordinator {
	return &WaitGroupCoordinator{}
}

// Start launches a registered service in a goroutine with proper coordination.
// This separates the concern of service lifecycle management from configuration.
func (c *WaitGroupCoordinator) Start(ctx context.Context, name string, fn func(ctx context.Context) error) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if err := fn(ctx); err != nil {
			log.Printf("[%s] service exited with error: %v", name, err)
		} else {
			log.Printf("[%s] service shutdown complete", name)
		}
	}()
}

// Shutdown waits for all registered services to complete.
// This follows the Interface Segregation Principle by providing a focused interface.
func (c *WaitGroupCoordinator) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All services shutdown complete")
		return nil
	case <-ctx.Done():
		log.Println("Shutdown timed out")
		return ctx.Err()
	}
}
