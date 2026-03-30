package pipeline

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RetryConfig struct {
	MaxAttempts int           // Default: 3
	InitialWait time.Duration // Default: 1s
	MaxWait     time.Duration // Default: 30s
	Multiplier  float64       // Default: 2.0
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		InitialWait: 1 * time.Second,
		MaxWait:     30 * time.Second,
		Multiplier:  2.0,
	}
}

func (c RetryConfig) withDefaults() RetryConfig {
	d := DefaultRetryConfig()
	if c.MaxAttempts <= 0 {
		c.MaxAttempts = d.MaxAttempts
	}
	if c.InitialWait <= 0 {
		c.InitialWait = d.InitialWait
	}
	if c.MaxWait <= 0 {
		c.MaxWait = d.MaxWait
	}
	if c.Multiplier <= 0 {
		c.Multiplier = d.Multiplier
	}
	return c
}

var (
	randFloat64        = defaultRandFloat64
	jitterDurationFn   = jitterDuration
	sleepWithContextFn = sleepWithContext
)

var (
	defaultRngMu sync.Mutex
	defaultRng   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func defaultRandFloat64() float64 {
	defaultRngMu.Lock()
	defer defaultRngMu.Unlock()
	return defaultRng.Float64()
}

func jitterDuration(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	factor := 0.9 + (randFloat64() * 0.2)
	return time.Duration(float64(d) * factor)
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}

	t := time.NewTimer(d)
	defer func() {
		if !t.Stop() {
			select {
			case <-t.C:
			default:
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func isRetryableUnwrapped(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if IsRetryable(e) {
			return true
		}
	}
	return false
}

func clampDuration(d, min, max time.Duration) time.Duration {
	if d < min {
		return min
	}
	if d > max {
		return max
	}
	return d
}

func Retry[T any](ctx context.Context, cfg RetryConfig, fn func() (T, error)) (T, error) {
	cfg = cfg.withDefaults()

	var zero T
	if ctx.Err() != nil {
		return zero, ctx.Err()
	}

	baseWait := min(cfg.InitialWait, cfg.MaxWait)

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return zero, ctx.Err()
		}

		fmt.Printf("Retry attempt %d/%d\n", attempt, cfg.MaxAttempts)

		val, err := fn()
		if err == nil {
			return val, nil
		}
		if ctx.Err() != nil {
			return zero, ctx.Err()
		}
		if !isRetryableUnwrapped(err) {
			return zero, err
		}
		if attempt == cfg.MaxAttempts {
			return zero, err
		}

		wait := clampDuration(baseWait, 0, cfg.MaxWait)
		wait = jitterDurationFn(wait)
		wait = clampDuration(wait, 0, cfg.MaxWait)
		if err := sleepWithContextFn(ctx, wait); err != nil {
			return zero, err
		}

		next := time.Duration(float64(baseWait) * cfg.Multiplier)
		baseWait = clampDuration(next, 0, cfg.MaxWait)
	}

	return zero, fmt.Errorf("pipeline.Retry: unreachable")
}

type FallbackChain[T any] struct {
	providers []func(context.Context) (T, error)
}

func NewFallbackChain[T any]() *FallbackChain[T] {
	return &FallbackChain[T]{}
}

func (fc *FallbackChain[T]) Add(provider func(context.Context) (T, error)) {
	if provider == nil {
		return
	}
	fc.providers = append(fc.providers, provider)
}

func (fc *FallbackChain[T]) Execute(ctx context.Context) (T, error) {
	var zero T
	if ctx.Err() != nil {
		return zero, ctx.Err()
	}
	if fc == nil || len(fc.providers) == 0 {
		return zero, fmt.Errorf("pipeline.FallbackChain.Execute: no providers")
	}

	var lastErr error
	for _, p := range fc.providers {
		if ctx.Err() != nil {
			return zero, ctx.Err()
		}
		v, err := p(ctx)
		if err == nil {
			return v, nil
		}
		lastErr = err
	}
	return zero, lastErr
}
