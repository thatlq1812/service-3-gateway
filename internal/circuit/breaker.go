package circuit

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State represents the circuit breaker state
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// Breaker implements circuit breaker pattern
type Breaker struct {
	maxFailures  uint32
	resetTimeout time.Duration

	mu              sync.RWMutex
	state           State
	failures        uint32
	lastFailTime    time.Time
	lastSuccessTime time.Time
}

// NewBreaker creates a new circuit breaker
func NewBreaker(maxFailures uint32, resetTimeout time.Duration) *Breaker {
	return &Breaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Execute runs the given function with circuit breaker protection
func (b *Breaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	// Check if circuit is open
	if b.isOpen() {
		return ErrCircuitOpen
	}

	// Execute the function
	err := fn(ctx)

	// Update circuit state based on result
	if err != nil {
		b.recordFailure()
		return err
	}

	b.recordSuccess()
	return nil
}

// isOpen checks if circuit breaker is open
func (b *Breaker) isOpen() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.state == StateClosed {
		return false
	}

	// Check if we should transition from Open to Half-Open
	if b.state == StateOpen && time.Since(b.lastFailTime) > b.resetTimeout {
		b.mu.RUnlock()
		b.mu.Lock()
		b.state = StateHalfOpen
		b.mu.Unlock()
		b.mu.RLock()
		return false
	}

	return b.state == StateOpen
}

// recordFailure records a failed execution
func (b *Breaker) recordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures++
	b.lastFailTime = time.Now()

	if b.state == StateHalfOpen {
		// Failed during half-open, immediately go back to open
		b.state = StateOpen
		return
	}

	// Check if we should open the circuit
	if b.failures >= b.maxFailures {
		b.state = StateOpen
	}
}

// recordSuccess records a successful execution
func (b *Breaker) recordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lastSuccessTime = time.Now()

	if b.state == StateHalfOpen {
		// Success during half-open, close the circuit
		b.state = StateClosed
		b.failures = 0
	}

	// Reset failure counter on success
	if b.state == StateClosed {
		b.failures = 0
	}
}

// GetState returns current circuit breaker state
func (b *Breaker) GetState() State {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.state
}

// GetStateString returns state as human-readable string
func (b *Breaker) GetStateString() string {
	state := b.GetState()
	switch state {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// GetFailures returns current failure count
func (b *Breaker) GetFailures() uint32 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.failures
}
