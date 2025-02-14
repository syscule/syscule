package lb

import (
	"fmt"
	"sync"
	"time"
)

// Dispatcher is a generic dispatcher that uses a load balancing strategy.
type Dispatcher struct {
	strategy LoadBalancer
	mu       sync.Mutex
}

// NewDispatcher creates a new Dispatcher with the given load balancing strategy.
func NewDispatcher(strategy LoadBalancer) *Dispatcher {
	return &Dispatcher{strategy: strategy}
}

// Dispatch uses the load balancing strategy to pick a target and then calls the provided dispatch function.
func (d *Dispatcher) Dispatch(dispatchFunc func(target *Target) error) error {
	d.mu.Lock()
	target := d.strategy.Pick()
	d.mu.Unlock()

	if target == nil {
		return fmt.Errorf("no available targets")
	}

	startTime := time.Now()
	target.IncrementActive()
	defer target.DecrementActive()

	err := dispatchFunc(target)
	duration := time.Since(startTime)
	target.UpdateResponseTime(duration)

	return err
}
