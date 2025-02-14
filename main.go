package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/syscule/syscule/pkg/lb"
)

const (
	CommandLoadBalancer       = "loadbalancer"
	StrategyLeastConnection   = "leastconnection"
	StrategyLeastResponseTime = "leastresponsetime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'loadbalancer' subcommands")
		return
	}

	command := os.Args[1]

	switch command {
	case CommandLoadBalancer:
		strategy := flag.String("strategy", "", "load balancing strategy")
		flag.CommandLine.Parse(os.Args[2:])

		if *strategy == "" {
			fmt.Println("expected strategy for 'loadbalancer' subcommand")
			return
		}

		switch *strategy {
		case StrategyLeastConnection:
			runLBLeastConnection()
		case StrategyLeastResponseTime:
			runLBLeastResponseTime()
		default:
			fmt.Printf("unknown strategy: %s\n", *strategy)
		}
	default:
		fmt.Printf("unknown command: %s\n", command)
	}
}

func runLBLeastConnection() {
	targets := []*lb.Target{
		{ID: "Server1", Active: 0},
		{ID: "Server2", Active: 0},
		{ID: "Server3", Active: 0},
	}

	lc := lb.NewLeastConnection(targets)
	dispatcher := lb.NewDispatcher(lc)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			err := dispatcher.Dispatch(func(target *lb.Target) error {
				target.IncrementActive()
				fmt.Printf("Request %d is being handled by %s\n", requestID, target.ID)
				time.Sleep(time.Millisecond * 100)
				target.DecrementActive()
				return nil
			})

			if err != nil {
				fmt.Printf("Request %d failed: %v\n", requestID, err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All requests have been processed")
}

func runLBLeastResponseTime() {
	targets := []*lb.Target{
		{ID: "Server1", ResponseTime: 100},
		{ID: "Server2", ResponseTime: 50},
		{ID: "Server3", ResponseTime: 150},
	}

	lrt := lb.NewLeastResponseTime(targets)
	dispatcher := lb.NewDispatcher(lrt)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			err := dispatcher.Dispatch(func(target *lb.Target) error {
				startTime := time.Now()
				target.IncrementActive()
				fmt.Printf("Request %d is being handled by %s\n", requestID, target.ID)
				time.Sleep(time.Millisecond * 100)
				duration := time.Since(startTime)
				target.UpdateResponseTime(duration)
				target.DecrementActive()
				return nil
			})

			if err != nil {
				fmt.Printf("Request %d failed: %v\n", requestID, err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All requests have been processed")
}
