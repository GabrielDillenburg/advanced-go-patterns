package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

type CircuitBreaker struct {
	state           State
	failureThreshold int
	resetTimeout     time.Duration
	failures         int
	lastFailure      time.Time
	mutex            sync.Mutex
}

func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
	}
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateClosed:
		return cb.executeClosed(operation)
	case StateHalfOpen:
		return cb.executeHalfOpen(operation)
	case StateOpen:
		return cb.executeOpen()
	default:
		return errors.New("unknown state")
	}
}

func (cb *CircuitBreaker) executeClosed(operation func() error) error {
	err := operation()
	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		if cb.failures >= cb.failureThreshold {
			cb.state = StateOpen
		}
		return err
	}
	cb.failures = 0
	return nil
}

func (cb *CircuitBreaker) executeHalfOpen(operation func() error) error {
	err := operation()
	if err != nil {
		cb.state = StateOpen
		cb.failures++
		cb.lastFailure = time.Now()
		return err
	}
	cb.state = StateClosed
	cb.failures = 0
	return nil
}

func (cb *CircuitBreaker) executeOpen() error {
	if time.Since(cb.lastFailure) > cb.resetTimeout {
		cb.state = StateHalfOpen
		return cb.executeHalfOpen(func() error {
			return errors.New("service is unavailable")
		})
	}
	return errors.New("circuit breaker is open")
}

func main() {
	cb := NewCircuitBreaker(3, 5*time.Second)

	// Simulating operations
	for i := 0; i < 10; i++ {
		err := cb.Execute(func() error {
			// Simulating a service call with random failures
			if rand.Float32() < 0.6 {
				return errors.New("service error")
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Operation %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("Operation %d succeeded\n", i+1)
		}

		time.Sleep(1 * time.Second)
	}
}