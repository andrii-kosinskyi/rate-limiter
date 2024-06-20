package ratelimiter

import (
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP := "192.168.1.1"

	// Тестирование разрешенных запросов
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Тестирование отклоненных запросов после превышения лимита
	if rl.Allow(clientIP) {
		t.Errorf("Expected request 6 to be denied")
	}

	// Ждем интервал для пополнения токенов
	time.Sleep(time.Second)

	// Тестирование разрешенных запросов после пополнения токенов
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed after refill", i+1)
		}
	}

	// Тестирование отклоненных запросов после повторного превышения лимита
	if rl.Allow(clientIP) {
		t.Errorf("Expected request 6 to be denied after refill")
	}
}

func TestRateLimiterDifferentIPs(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP1 := "192.168.1.1"
	clientIP2 := "192.168.1.2"

	// Тестирование разрешенных запросов для первого IP
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP1) {
			t.Errorf("Expected request %d to be allowed for IP1", i+1)
		}
	}

	// Тестирование отклоненных запросов для первого IP после превышения лимита
	if rl.Allow(clientIP1) {
		t.Errorf("Expected request 6 to be denied for IP1")
	}

	// Тестирование разрешенных запросов для второго IP
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP2) {
			t.Errorf("Expected request %d to be allowed for IP2", i+1)
		}
	}

	// Тестирование отклоненных запросов для второго IP после превышения лимита
	if rl.Allow(clientIP2) {
		t.Errorf("Expected request 6 to be denied for IP2")
	}
}

func TestRateLimiterPartialRefill(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP := "192.168.1.1"

	// Используем все токены
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Ждем полсекунды для частичного пополнения токенов
	time.Sleep(500 * time.Millisecond)

	// Проверяем, что после частичного пополнения доступен один токен
	if !rl.Allow(clientIP) {
		t.Errorf("Expected request after partial refill to be allowed")
	}

	// Следующий запрос должен быть отклонен, так как токенов больше нет
	if rl.Allow(clientIP) {
		t.Errorf("Expected request after partial refill to be denied")
	}
}

func TestRateLimiterMultipleRefills(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)

	clientIP := "192.168.1.1"

	// Используем все токены
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Ждем два интервала для полного пополнения токенов
	time.Sleep(2 * time.Second)

	// Проверяем, что после полного пополнения доступно пять токенов
	for i := 0; i < 5; i++ {
		if !rl.Allow(clientIP) {
			t.Errorf("Expected request %d to be allowed after full refill", i+1)
		}
	}

	// Следующий запрос должен быть отклонен, так как токенов больше нет
	if rl.Allow(clientIP) {
		t.Errorf("Expected request 6 to be denied after full refill")
	}
}
