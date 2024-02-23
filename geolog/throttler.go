package main

import (
	"sync"
	"time"
)

type ipThrottler struct {
	mu  sync.Mutex
	ips map[string]time.Time
}

func NewIPThrottler() *ipThrottler {
	return &ipThrottler{
		ips: make(map[string]time.Time),
	}
}

func (t *ipThrottler) Allow(ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if lastTime, ok := t.ips[ip]; ok && time.Since(lastTime) < 5*time.Second {
		return false
	}

	t.ips[ip] = time.Now()
	return true
}
