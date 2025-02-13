package lb_test

import (
	"sync"
	"testing"

	"github.com/syscule/syscule/pkg/lb"
)

func TestLeastConnection_Pick(t *testing.T) {
	tests := []struct {
		name     string
		targets  []*lb.Target
		expected string
	}{
		{
			name: "Single target",
			targets: []*lb.Target{
				{ID: "A", Active: 0},
			},
			expected: "A",
		},
		{
			name: "Multiple targets, one least loaded",
			targets: []*lb.Target{
				{ID: "A", Active: 2},
				{ID: "B", Active: 1},
				{ID: "C", Active: 3},
			},
			expected: "B",
		},
		{
			name: "Multiple targets, tie",
			targets: []*lb.Target{
				{ID: "A", Active: 2},
				{ID: "B", Active: 2},
				{ID: "C", Active: 3},
			},
			expected: "A",
		},
		{
			name:     "No targets",
			targets:  []*lb.Target{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := lb.NewLeastConnection(tt.targets)
			target := lc.Pick()

			if target == nil {
				if tt.expected != "" {
					t.Errorf("expected %s, got nil", tt.expected)
				}
			} else if target.ID != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, target.ID)
			}
		})
	}
}

func TestLeastConnection_Concurrency(t *testing.T) {
	targets := []*lb.Target{
		{ID: "A", Active: 0},
		{ID: "B", Active: 0},
		{ID: "C", Active: 0},
	}

	lc := lb.NewLeastConnection(targets)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			target := lc.Pick()
			if target != nil {
				target.Increment()
				target.Decrement()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
