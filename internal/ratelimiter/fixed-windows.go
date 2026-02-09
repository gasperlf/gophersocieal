package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	sync.Mutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	count, ok := rl.clients[ip]
	if !ok {
		rl.clients[ip] = 1
		go rl.resetCount(ip)
		return true, 0
	}

	if count < rl.limit {
		rl.clients[ip]++
		return true, 0
	}

	return false, 0
}

func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	defer rl.Unlock()
	delete(rl.clients, ip)
}
