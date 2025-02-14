package lb_test

import (
	"fmt"
	"testing"

	"github.com/syscule/syscule/pkg/lb"
)

func TestDispatcher_LeastConnection(t *testing.T) {
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
			dispatcher := lb.NewDispatcher(lc)

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

func TestDispatcher_LeastResponseTime(t *testing.T) {
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
