package lb

import (
	"sync"
)

type LeastResponseTime struct {
	targets []*Target
	mu      sync.Mutex
}

func NewLeastResponseTime(targets []*Target) *LeastResponseTime {
	return &LeastResponseTime{targets: targets}
}

func (lrt *LeastResponseTime) Pick() *Target {
	lrt.mu.Lock()
	defer lrt.mu.Unlock()

	if len(lrt.targets) == 0 {
		return nil
	}

	var leastLoaded *Target
	minResponseTime := int(^uint(0) >> 1)

	for _, target := range lrt.targets {
		target.mu.Lock()
		responseTime := target.ResponseTime
		target.mu.Unlock()

		if responseTime < minResponseTime {
			leastLoaded = target
			minResponseTime = responseTime
		}
	}

	return leastLoaded
}

func (t *Target) UpdateResponseTime(responseTime int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ResponseTime = responseTime
}
