package lb

import (
	"sync"
	"time"
)

// Target is a target for load balancing.
type Target struct {
	ID           string
	Active       int
	ResponseTime int
	mu           sync.Mutex
}

func (t *Target) IncrementActive() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Active++
}

func (t *Target) DecrementActive() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Active--
}

func (t *Target) UpdateResponseTime(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ResponseTime = int(duration.Milliseconds())
}
