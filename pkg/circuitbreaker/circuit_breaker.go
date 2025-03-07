package use_CircuitBreaker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type CircuitBreaker interface {
	Execute(ctx context.Context, command func() error) error
	State() State
	Reset()
}
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

type Options struct {
	Timeout          time.Duration
	MaxFailures      int
	ResetTimeout     time.Duration
	HalfOpenMaxCalls int
}

func NewOptions() *Options {
	return &Options{
		Timeout:          5 * time.Second,
		MaxFailures:      5,
		ResetTimeout:     1 * time.Minute,
		HalfOpenMaxCalls: 3,
	}
}

type DefaultCircuitBreaker struct {
	state       atomic.Value
	options     *Options
	failures    atomic.Int32
	lastFailure atomic.Value
	calls       atomic.Int32
}

func NewCircuitBreaker(options *Options) CircuitBreaker {
	if options == nil {
		options = NewOptions()
	}
	cb := &DefaultCircuitBreaker{
		options: options,
	}
	cb.state.Store(StateClosed)
	cb.lastFailure.Store(time.Time{})
	return cb
}
func (cb *DefaultCircuitBreaker) Execute(ctx context.Context, command func() error) error {
	state := cb.state.Load().(State)
	switch state {
	case StateOpen:
		lastFailure := cb.lastFailure.Load().(time.Time)
		if time.Since(lastFailure) > cb.options.ResetTimeout {
			cb.state.Store(StateHalfOpen)
			cb.calls.Store(0)
		} else {
			return fmt.Errorf("circuit breaker is open: waiting for reset timeout (%v)", cb.options.ResetTimeout)
		}
	case StateHalfOpen:
		if cb.calls.Load() >= int32(cb.options.HalfOpenMaxCalls) {
			return fmt.Errorf("too many calls in half-open state (max: %d)", cb.options.HalfOpenMaxCalls)
		}
	}
	cb.calls.Add(1)
	select {
	case <-ctx.Done():
		return fmt.Errorf("command execution cancelled: %w", ctx.Err())
	default:
		err := command()
		if err != nil {
			failures := cb.failures.Add(1)
			cb.lastFailure.Store(time.Now())
			if int(failures) >= cb.options.MaxFailures {
				cb.state.Store(StateOpen)
			}
			return fmt.Errorf("command execution failed: %w", err)
		}
		if cb.state.Load() == StateHalfOpen {
			cb.state.Store(StateClosed)
		}
		cb.failures.Store(0)
		return nil
	}
}
func (cb *DefaultCircuitBreaker) State() State {
	return cb.state.Load().(State)
}
func (cb *DefaultCircuitBreaker) Reset() {
	cb.state.Store(StateClosed)
	cb.failures.Store(0)
	cb.calls.Store(0)
	cb.lastFailure.Store(time.Time{})
}
