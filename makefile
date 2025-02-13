.PHONY: all build run clean

# Binary name
BINARY := loadbalancer

# Default target
all: build

# Build
build:
	go build -o $(BINARY) main.go

# Load Balancers
## Least Connection
run-lb-lc: build
	./$(BINARY) loadbalancer -strategy=leastconnection

## Least Response Time
run-lb-lrt: build
	./$(BINARY) loadbalancer -strategy=leastresponsetime

# Clean up
clean:
	rm -f $(BINARY)