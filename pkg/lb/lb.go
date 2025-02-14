package lb

// LoadBalancer is an interface for load balancing strategies.
type LoadBalancer interface {
	Pick() *Target
	Calculate(target *Target) int
}
