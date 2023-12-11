package healthcheck

import "time"

type HealthCheck func(timeout time.Duration) Result

type Result struct {
	Message  string
	Status   Status
	Duration time.Duration
	Err      error
}

type Status int

const (
	Green  Status = iota + 1 // Healthy
	Yellow                   // Healthy but requires attention
	Red                      // Not healthy
)
