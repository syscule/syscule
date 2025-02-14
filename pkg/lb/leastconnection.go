package lb

import (
	"sync"
)

// LeastConnection implements the LoadBalancer interface using the least connection strategy.
type LeastConnection struct {
	targets []*Target
	mu      sync.Mutex
}

// NewLeastConnection creates a new LeastConnection load balancer.
func NewLeastConnection(targets []*Target) *LeastConnection {
	return &LeastConnection{targets: targets}
}

// Pick selects the target with the least active connections.
func (lc *LeastConnection) Pick() *Target {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if len(lc.targets) == 0 {
		return nil
	}

	var bestTarget *Target
	for _, target := range lc.targets {
		if bestTarget == nil || lc.Calculate(target) < lc.Calculate(bestTarget) {
			bestTarget = target
		}
	}

	return bestTarget
}

// Calculate returns the number of active connections of the target.
func (lc *LeastConnection) Calculate(target *Target) int {
	target.mu.Lock()
	defer target.mu.Unlock()
	return target.Active
}
