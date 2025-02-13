package lb_test

import (
	"sync"
	"testing"

	"github.com/syscule/syscule/pkg/lb"
)

func TestLeastResponseTime_Pick(t *testing.T) {
	tests := []struct {
		name     string
		targets  []*lb.Target
		expected string
	}{
		{
			name: "Single target",
			targets: []*lb.Target{
				{ID: "A", ResponseTime: 100},
			},
			expected: "A",
		},
		{
			name: "Multiple targets, one least response time",
			targets: []*lb.Target{
				{ID: "A", ResponseTime: 200},
				{ID: "B", ResponseTime: 150},
				{ID: "C", ResponseTime: 250},
			},
			expected: "B",
		},
		{
			name: "Multiple targets, tie",
			targets: []*lb.Target{
				{ID: "A", ResponseTime: 150},
				{ID: "B", ResponseTime: 150},
				{ID: "C", ResponseTime: 200},
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
			lrt := lb.NewLeastResponseTime(tt.targets)
			target := lrt.Pick()

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

func TestLeastResponseTime_Concurrency(t *testing.T) {
	targets := []*lb.Target{
		{ID: "A", ResponseTime: 100},
		{ID: "B", ResponseTime: 120},
		{ID: "C", ResponseTime: 110},
	}

	lrt := lb.NewLeastResponseTime(targets)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			target := lrt.Pick()
			if target != nil {
				target.UpdateResponseTime(target.ResponseTime + 10) // Increasing load simulation
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
