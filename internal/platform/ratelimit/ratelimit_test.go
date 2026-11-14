package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiterFixedWindow(t *testing.T) {
	now := time.Date(2026, 11, 14, 10, 0, 0, 0, time.UTC)
	limiter := New(2, time.Minute)
	limiter.now = func() time.Time { return now }

	if ok, _ := limiter.Allow("ip"); !ok {
		t.Fatal("first attempt denied")
	}
	if ok, _ := limiter.Allow("ip"); !ok {
		t.Fatal("second attempt denied")
	}
	ok, retry := limiter.Allow("ip")
	if ok || retry <= 0 {
		t.Fatalf("third attempt should be denied with retry: ok=%v retry=%v", ok, retry)
	}
	if ok, _ := limiter.Allow("other"); !ok {
		t.Fatal("independent key denied")
	}

	now = now.Add(2 * time.Minute)
	if ok, _ := limiter.Allow("ip"); !ok {
		t.Fatal("attempt after window denied")
	}
}

func TestLimiterStaysBoundedAndEvictsEarliest(t *testing.T) {
	now := time.Date(2026, 11, 14, 10, 0, 0, 0, time.UTC)
	limiter := New(1, time.Minute)
	limiter.now = func() time.Time { return now }
	for i := 0; i < maxKeys; i++ {
		key := fmt.Sprintf("key-%d", i)
		if ok, _ := limiter.Allow(key); !ok {
			t.Fatalf("attempt %d denied", i)
		}
		now = now.Add(time.Millisecond)
	}
	if len(limiter.items) != maxKeys {
		t.Fatalf("limiter grew to %d", len(limiter.items))
	}
	if ok, _ := limiter.Allow("newcomer"); !ok {
		t.Fatal("new caller denied")
	}
	if len(limiter.items) != maxKeys {
		t.Fatalf("limiter not bounded: %d", len(limiter.items))
	}
	if _, found := limiter.items["key-0"]; found {
		t.Fatal("oldest entry was not evicted")
	}
	if _, found := limiter.items["newcomer"]; !found {
		t.Fatal("new entry missing")
	}
}

func TestLimiterResetAndDisabled(t *testing.T) {
	limiter := New(1, time.Minute)
	if ok, _ := limiter.Allow("ip"); !ok {
		t.Fatal("first attempt denied")
	}
	if ok, _ := limiter.Allow("ip"); ok {
		t.Fatal("second attempt should be denied")
	}
	limiter.Reset("ip")
	if ok, _ := limiter.Allow("ip"); !ok {
		t.Fatal("attempt after reset denied")
	}

	var disabled *Limiter
	if ok, _ := disabled.Allow("ip"); !ok {
		t.Fatal("nil limiter should allow")
	}
	disabled.Reset("ip")
	if ok, _ := New(0, time.Minute).Allow("ip"); !ok {
		t.Fatal("zero max should allow")
	}
}
