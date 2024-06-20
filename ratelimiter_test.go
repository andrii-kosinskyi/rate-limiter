package ratelimiter

import (
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP := "192.168.1.1"

	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Testing denied requests after exceeding the limit
	if rl.Allow(clientIP) {
		t.Errorf("Expected request 6 to be denied")
	}

	// Waiting for the interval to refill tokens
	time.Sleep(time.Second)

	// Testing allowed requests after token refill
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed after refill", i+1)
		}
	}

	// Testing denied requests after exceeding the limit again
	if rl.Allow(clientIP) {
		t.Errorf("Expected request 6 to be denied after refill")
	}
}

func TestRateLimiterDifferentIPs(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP1 := "192.168.1.1"
	clientIP2 := "192.168.1.2"

	// Testing allowed requests for the first IP
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP1) {
			t.Errorf("Expected request %d to be allowed for IP1", i+1)
		}
	}

	// Testing denied requests for the first IP after exceeding the limit
	if rl.Allow(clientIP1) {
		t.Errorf("Expected request 6 to be denied for IP1")
	}

	// Testing allowed requests for the second IP
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP2) {
			t.Errorf("Expected request %d to be allowed for IP2", i+1)
		}
	}

	// Testing denied requests for the second IP after exceeding the limit
	if rl.Allow(clientIP2) {
		t.Errorf("Expected request 6 to be denied for IP2")
	}
}

func TestRateLimiterPartialRefill(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP := "192.168.1.1"

	// Using up all the tokens
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Waiting for quarter a second for partial token refill
	time.Sleep(200 * time.Millisecond)

	// Checking that one token is available after partial refill
	if !rl.Allow(clientIP) {
		t.Errorf("Expected request after partial refill to be allowed")
	}

	// The next request should be denied as there are no more tokens
	if rl.Allow(clientIP) {
		t.Errorf("Expected request after partial refill to be denied")
	}
}

func TestRateLimiterMultipleRefills(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP := "192.168.1.1"

	// Using up all the tokens
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Waiting for two intervals for full token refill
	time.Sleep(2 * time.Second)

	// Checking that five tokens are available after full refill
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed after full refill", i+1)
		}
	}

	// The next request should be denied as there are no more tokens
	if rl.Allow(clientIP) {
		t.Errorf("Expected request 6 to be denied after full refill")
	}
}
