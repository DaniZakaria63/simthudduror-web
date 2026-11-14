// Package ratelimit provides the small in-process limiter used for
// authentication attempts. It is intentionally bounded and simple.
package ratelimit

import (
	"sync"
	"time"
)

const maxKeys = 4096

type entry struct {
	count   int
	resetAt time.Time
}

// Limiter is a fixed-window counter keyed by caller. It is safe for concurrent
// use. The map is capped at maxKeys; expired entries are swept periodically and
// the earliest-expiring entry is evicted when a new caller arrives at capacity.
type Limiter struct {
	mu        sync.Mutex
	max       int
	window    time.Duration
	now       func() time.Time
	nextSweep time.Time
	items     map[string]entry
}

// New creates a limiter that permits max attempts per window.
func New(max int, window time.Duration) *Limiter {
	return &Limiter{max: max, window: window, now: time.Now, items: make(map[string]entry)}
}

// Allow reports whether key may proceed and, when denied, how long remains.
func (l *Limiter) Allow(key string) (bool, time.Duration) {
	if l == nil || l.max <= 0 || l.window <= 0 {
		return true, 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	l.sweepLocked(now)

	item, found := l.items[key]
	if !found {
		if len(l.items) >= maxKeys {
			l.evictEarliestLocked()
		}
		l.items[key] = entry{count: 1, resetAt: now.Add(l.window)}
		return true, 0
	}
	if !now.Before(item.resetAt) {
		l.items[key] = entry{count: 1, resetAt: now.Add(l.window)}
		return true, 0
	}
	if item.count >= l.max {
		return false, item.resetAt.Sub(now)
	}
	item.count++
	l.items[key] = item
	return true, 0
}

// Reset clears one caller, used after a successful authentication.
func (l *Limiter) Reset(key string) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.items, key)
}

// sweepLocked removes expired entries at most once per window rather than on
// every request.
func (l *Limiter) sweepLocked(now time.Time) {
	if now.Before(l.nextSweep) {
		return
	}
	for key, item := range l.items {
		if !now.Before(item.resetAt) {
			delete(l.items, key)
		}
	}
	l.nextSweep = now.Add(l.window)
}

func (l *Limiter) evictEarliestLocked() {
	var oldestKey string
	var oldest time.Time
	for key, item := range l.items {
		if oldestKey == "" || item.resetAt.Before(oldest) {
			oldestKey, oldest = key, item.resetAt
		}
	}
	if oldestKey != "" {
		delete(l.items, oldestKey)
	}
}
