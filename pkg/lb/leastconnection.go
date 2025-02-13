package lb

import (
	"sync"
)

type Target struct {
	ID     string
	Active int
	mu     sync.Mutex
}

func (t *Target) Increment() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Active++
}

func (t *Target) Decrement() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Active--
}

type LeastConnection struct {
	targets []*Target
	mu      sync.Mutex
}

func NewLeastConnection(targets []*Target) *LeastConnection {
	return &LeastConnection{targets: targets}
}

func (lc *LeastConnection) Pick() *Target {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if len(lc.targets) == 0 {
		return nil
	}

	var leastLoaded *Target
	for _, target := range lc.targets {
		if leastLoaded == nil || target.Active < leastLoaded.Active {
			leastLoaded = target
		}
	}

	return leastLoaded
}
