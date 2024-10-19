package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type CircuitBreaker struct {
	mu             sync.Mutex
	failures       int
	failureLimit   int
	open           bool
	openUntil      time.Time
	retryTimeout   time.Duration
	consecutiveReq int
}

// NewCircuitBreaker creates a new CircuitBreaker with a failure limit and retry timeout
func NewCircuitBreaker(failureLimit int, retryTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureLimit: failureLimit,
		retryTimeout: retryTimeout,
	}
}

// Call executes the service, handling the circuit breaker logic
func (cb *CircuitBreaker) Call(service func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.open {
		// If the circuit is open and the retry timeout has not elapsed, fail fast
		if time.Now().Before(cb.openUntil) {
			return errors.New("circuit breaker is open")
		}
		// Half-open state, allow one request
		cb.open = false
		cb.failures = 0
	}

	// Call the service
	err := service()
	if err != nil {
		cb.failures++
		if cb.failures >= cb.failureLimit {
			// Open the circuit breaker if failure limit is reached
			cb.open = true
			cb.openUntil = time.Now().Add(cb.retryTimeout)
		}
		return err
	}

	// Reset failure count on success
	cb.failures = 0
	return nil
}

// Example service function that may fail
func exampleService() error {
	// Simulate a failure scenario
	return errors.New("service failure")
}

func main() {
	// Initialize a new circuit breaker with a failure limit of 3 and a retry timeout of 10 seconds
	cb := NewCircuitBreaker(3, 10*time.Second)

	for i := 0; i < 10; i++ {
		err := cb.Call(exampleService)
		if err != nil {
			fmt.Printf("Request %d failed: %s\n", i+1, err)
		} else {
			fmt.Printf("Request %d succeeded\n", i+1)
		}

		// Simulate a delay between requests
		time.Sleep(1 * time.Second)
	}
}
