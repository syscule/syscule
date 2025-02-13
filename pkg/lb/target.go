package lb

import (
	"sync"
)

type Target struct {
	ID           string
	Active       int
	ResponseTime int
	mu           sync.Mutex
}
