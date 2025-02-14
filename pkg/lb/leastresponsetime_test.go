package lb_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

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
			dispatcher := lb.NewDispatcher(lrt)

			err := dispatcher.Dispatch(func(target *lb.Target) error {
				if target.ID != tt.expected {
					return fmt.Errorf("expected %s, got %s", tt.expected, target.ID)
				}
				return nil
			})

			if err != nil && tt.expected != "" {
				t.Errorf("dispatch failed: %v", err)
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
	dispatcher := lb.NewDispatcher(lrt)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			err := dispatcher.Dispatch(func(target *lb.Target) error {
				startTime := time.Now()
				target.IncrementActive()
				defer target.DecrementActive()
				time.Sleep(time.Millisecond * 10) // Simulate work
				duration := time.Since(startTime)
				target.UpdateResponseTime(duration)
				return nil
			})
			if err != nil {
				t.Errorf("dispatch failed: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
