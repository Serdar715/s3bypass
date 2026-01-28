package limiter

import (
	"math/rand"
	"s3bypass/pkg/config"
	"time"
)

// RateLimiter handles request delays and jitter
type RateLimiter struct {
	BaseDelay int // Milliseconds
}

// New creates a new RateLimiter
func New(baseDelay int) *RateLimiter {
	return &RateLimiter{
		BaseDelay: baseDelay,
	}
}

// Wait blocks for the calculated delay duration
func (l *RateLimiter) Wait() {
	if l.BaseDelay <= 0 {
		return
	}

	// Calculate jitter (+/- 10%)
	jitter := int(float64(l.BaseDelay) * config.JitterPercentage)
	
	// Formula: delay + rand(2*jitter + 1) - jitter
	// This creates a range of [delay-jitter, delay+jitter]
	randomPart := rand.Intn(jitter*config.JitterMultiplier + config.JitterBaseOffset)
	actualDelay := l.BaseDelay + randomPart - jitter

	if actualDelay < 0 {
		actualDelay = 0
	}

	time.Sleep(time.Duration(actualDelay) * time.Millisecond)
}
