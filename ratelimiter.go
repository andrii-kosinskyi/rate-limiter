package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate     int
	interval time.Duration
	buckets  map[string]*tokenBucket
	mu       sync.Mutex
}

type tokenBucket struct {
	tokens     int
	lastRefill time.Time
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		interval: interval,
		buckets:  make(map[string]*tokenBucket),
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[ip]
	if !exists {
		bucket = &tokenBucket{
			tokens:     rl.rate,
			lastRefill: time.Now(),
		}
		rl.buckets[ip] = bucket
	}
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	if elapsed >= rl.interval {
		bucket.tokens = rl.rate
		bucket.lastRefill = now
	} else {
		refillTokens := int(elapsed.Seconds() / rl.interval.Seconds() * float64(rl.rate))
		bucket.tokens = min(bucket.tokens+refillTokens, rl.rate)
		bucket.lastRefill = bucket.lastRefill.Add(time.Duration(refillTokens) * rl.interval / time.Duration(rl.rate))
	}
	if bucket.tokens > 0 {
		bucket.tokens--

		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
