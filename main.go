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
		os.Exit(1)
	}

	switch os.Args[1] {
	case CommandLoadBalancer:
		loadBalancerCmd := flag.NewFlagSet(CommandLoadBalancer, flag.ExitOnError)
		strategy := loadBalancerCmd.String("strategy", StrategyLeastConnection, "Load balancing strategy")
		err := loadBalancerCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing command:", err)
			os.Exit(1)
		}

		switch *strategy {
		case StrategyLeastConnection:
			runLBLeastConnection()
		case StrategyLeastResponseTime:
			runLBLeastResponseTime()
		default:
			fmt.Println("Unknown strategy:", *strategy)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func runLBLeastConnection() {
	targets := []*lb.Target{
		{ID: "Server1", Active: 0},
		{ID: "Server2", Active: 0},
		{ID: "Server3", Active: 0},
	}

	lc := lb.NewLeastConnection(targets)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()
			target := lc.Pick()
			if target != nil {
				target.Increment()
				fmt.Printf("Request %d is being handled by %s\n", requestID, target.ID)
				time.Sleep(time.Millisecond * 100)
				target.Decrement()
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

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			target := lrt.Pick()
			if target != nil {
				target.Increment()
				fmt.Printf("Request %d is being handled by %s\n", requestID, target.ID)

				time.Sleep(time.Millisecond * 100)
				target.UpdateResponseTime(target.ResponseTime + 5)
				target.Decrement()
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All requests have been processed")
}
