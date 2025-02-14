package lb

import (
	"sync"
)

// LeastResponseTime implements the LoadBalancer interface using the least response time strategy.
type LeastResponseTime struct {
	targets []*Target
	mu      sync.Mutex
}

// NewLeastResponseTime creates a new LeastResponseTime load balancer.
func NewLeastResponseTime(targets []*Target) *LeastResponseTime {
	return &LeastResponseTime{targets: targets}
}

// Pick selects the target with the least response time.
func (lrt *LeastResponseTime) Pick() *Target {
	lrt.mu.Lock()
	defer lrt.mu.Unlock()

	if len(lrt.targets) == 0 {
		return nil
	}

	var bestTarget *Target
	for _, target := range lrt.targets {
		if bestTarget == nil || lrt.Calculate(target) < lrt.Calculate(bestTarget) {
			bestTarget = target
		}
	}

	return bestTarget
}

// Calculate returns the response time of the target.
func (lrt *LeastResponseTime) Calculate(target *Target) int {
	target.mu.Lock()
	defer target.mu.Unlock()
	return target.ResponseTime
}
