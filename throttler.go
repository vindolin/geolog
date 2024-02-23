package main

import (
	"sync"
	"time"
)

// IPThrottler throttles requests based on IP address.
type IPThrottler struct {
	mu        sync.Mutex
	ips       map[string]time.Time
	throttle  time.Duration
	cleanup   time.Duration
	lastClean time.Time
}

// New creates a new IPThrottler.
func NewIPThrottler(throttle, cleanup time.Duration) *IPThrottler {
	return &IPThrottler{
		ips:       make(map[string]time.Time),
		throttle:  throttle,
		cleanup:   cleanup,
		lastClean: time.Now(),
	}
}

// Allow checks if a request from the given IP is allowed.
func (t *IPThrottler) Allow(ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Clean up old entries if necessary
	if time.Since(t.lastClean) > t.cleanup {
		for ip, lastTime := range t.ips {
			if time.Since(lastTime) > t.throttle {
				delete(t.ips, ip)
			}
		}
		t.lastClean = time.Now()
	}

	// Check if request is allowed
	if lastTime, ok := t.ips[ip]; ok && time.Since(lastTime) < t.throttle {
		return false
	}

	t.ips[ip] = time.Now()
	return true
}
