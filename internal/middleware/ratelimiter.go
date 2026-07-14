package middleware

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
}

const (
	// Maximum sustained requests per second
	rateLimit = 10

	// Maximum burst allowed
	burstSize = 5

	// Cleanup interval for inactive visitors
	cleanupEvery = time.Minute

	// Remove visitors inactive for this duration
	visitorTimeout = 10 * time.Minute
)

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
	}

	go rl.cleanupVisitors()
	return rl

}

func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(cleanupEvery)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()

		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > visitorTimeout {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	now := time.Now()
	if !exists {
		visitor = &Visitor{
			limiter:  rate.NewLimiter(rate.Every(time.Second/rateLimit), burstSize),
			lastSeen: now,
		}
		rl.visitors[ip] = visitor
	}

	visitor.lastSeen = now

	return visitor
}
