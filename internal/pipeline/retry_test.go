package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry_ExponentialBackoff(t *testing.T) {
	oldJitter := jitterDurationFn
	oldSleep := sleepWithContextFn
	defer func() {
		jitterDurationFn = oldJitter
		sleepWithContextFn = oldSleep
	}()

	jitterDurationFn = func(d time.Duration) time.Duration { return d }
	var sleeps []time.Duration
	sleepWithContextFn = func(ctx context.Context, d time.Duration) error {
		sleeps = append(sleeps, d)
		return nil
	}

	cfg := RetryConfig{MaxAttempts: 4, InitialWait: 1 * time.Second, MaxWait: 30 * time.Second, Multiplier: 2.0}
	_, err := Retry(context.Background(), cfg, func() (int, error) {
		return 0, ErrTrnTimeoutErr("x", errors.New("timeout"))
	})
	if err == nil {
		t.Fatalf("expected error")
	}

	if len(sleeps) != 3 {
		t.Fatalf("expected 3 sleeps, got %d", len(sleeps))
	}
	if sleeps[0] != 1*time.Second || sleeps[1] != 2*time.Second || sleeps[2] != 4*time.Second {
		t.Fatalf("unexpected backoff sleeps: %v", sleeps)
	}
}

func TestRetry_NonRetryable(t *testing.T) {
	oldJitter := jitterDurationFn
	oldSleep := sleepWithContextFn
	defer func() {
		jitterDurationFn = oldJitter
		sleepWithContextFn = oldSleep
	}()

	jitterDurationFn = func(d time.Duration) time.Duration { return d }
	calledSleep := 0
	sleepWithContextFn = func(ctx context.Context, d time.Duration) error {
		calledSleep++
		return nil
	}

	calledFn := 0
	cfg := RetryConfig{MaxAttempts: 5, InitialWait: 10 * time.Millisecond, MaxWait: 1 * time.Second, Multiplier: 2.0}
	_, err := Retry(context.Background(), cfg, func() (int, error) {
		calledFn++
		return 0, ErrTrnAPIKeyErr("x")
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if calledFn != 1 {
		t.Fatalf("expected fn called once, got %d", calledFn)
	}
	if calledSleep != 0 {
		t.Fatalf("expected no sleeps, got %d", calledSleep)
	}
}

func TestRetry_Success(t *testing.T) {
	oldJitter := jitterDurationFn
	oldSleep := sleepWithContextFn
	defer func() {
		jitterDurationFn = oldJitter
		sleepWithContextFn = oldSleep
	}()

	jitterDurationFn = func(d time.Duration) time.Duration { return d }
	calledSleep := 0
	sleepWithContextFn = func(ctx context.Context, d time.Duration) error {
		calledSleep++
		return nil
	}

	attempts := 0
	cfg := RetryConfig{MaxAttempts: 5, InitialWait: 1 * time.Millisecond, MaxWait: 1 * time.Second, Multiplier: 2.0}
	v, err := Retry(context.Background(), cfg, func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", ErrRwtTimeoutErr("x", errors.New("timeout"))
		}
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if v != "ok" {
		t.Fatalf("expected ok, got %q", v)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	if calledSleep != 2 {
		t.Fatalf("expected 2 sleeps, got %d", calledSleep)
	}
}

func TestRetry_ContextCancel(t *testing.T) {
	oldJitter := jitterDurationFn
	oldSleep := sleepWithContextFn
	defer func() {
		jitterDurationFn = oldJitter
		sleepWithContextFn = oldSleep
	}()

	jitterDurationFn = func(d time.Duration) time.Duration { return d }
	sleepWithContextFn = func(ctx context.Context, d time.Duration) error {
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	attempts := 0
	_, err := Retry(ctx, RetryConfig{MaxAttempts: 5, InitialWait: 10 * time.Second, MaxWait: 30 * time.Second, Multiplier: 2.0}, func() (int, error) {
		attempts++
		if attempts == 1 {
			cancel()
			return 0, ErrTrnTimeoutErr("x", errors.New("timeout"))
		}
		return 1, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got: %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestFallbackChain(t *testing.T) {
	fc := NewFallbackChain[int]()
	order := []int{}

	ctx := context.WithValue(context.Background(), "k", "v")

	fc.Add(func(ctx context.Context) (int, error) {
		if ctx.Value("k") != "v" {
			return 0, errors.New("missing ctx value")
		}
		order = append(order, 1)
		return 0, errors.New("fail-1")
	})
	fc.Add(func(ctx context.Context) (int, error) {
		order = append(order, 2)
		return 42, nil
	})
	fc.Add(func(ctx context.Context) (int, error) {
		order = append(order, 3)
		return 0, errors.New("should not be called")
	})

	v, err := fc.Execute(ctx)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("unexpected provider order: %v", order)
	}
}
