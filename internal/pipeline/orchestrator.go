package pipeline

import (
	"context"
	"errors"
	"fmt"
)

type StepFunc func(ctx context.Context, projectID string) (State, error)

type EventEmitter interface {
	Emit(event string, data interface{})
}

type OrchestratorConfig struct {
	Emitter     EventEmitter
	RetryConfig RetryConfig
	StateStore  StateStore
}

type StateStore interface {
	GetState(projectID string) (State, error)
	SetState(projectID string, state State) error
}

type Orchestrator struct {
	config OrchestratorConfig
	steps  map[State]StepFunc
}

func NewOrchestrator(cfg OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		config: cfg,
		steps:  make(map[State]StepFunc),
	}
}
func (o *Orchestrator) RegisterStep(state State, fn StepFunc) {
	if o == nil {
		return
	}
	if o.steps == nil {
		o.steps = make(map[State]StepFunc)
	}
	o.steps[state] = fn
}
func (o *Orchestrator) Run(ctx context.Context, projectID string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if o == nil {
		return fmt.Errorf("pipeline.Orchestrator.Run: nil orchestrator")
	}
	if o.config.StateStore == nil {
		return fmt.Errorf("pipeline.Orchestrator.Run: nil StateStore")
	}

	st, err := o.config.StateStore.GetState(projectID)
	if err != nil {
		return fmt.Errorf("pipeline.Orchestrator.Run: get state: %w", err)
	}
	if st == "" {
		st = StateOrder[0]
	}
	return o.runFromState(ctx, projectID, st)
}

func (o *Orchestrator) RunFrom(ctx context.Context, projectID string, fromState State) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if o == nil {
		return fmt.Errorf("pipeline.Orchestrator.RunFrom: nil orchestrator")
	}
	if o.config.StateStore == nil {
		return fmt.Errorf("pipeline.Orchestrator.RunFrom: nil StateStore")
	}
	if err := o.config.StateStore.SetState(projectID, fromState); err != nil {
		wrapped := fmt.Errorf("pipeline.Orchestrator.RunFrom: set state %s: %w", fromState, err)
		o.emit("pipeline:error", map[string]interface{}{
			"projectID": projectID,
			"step":      string(fromState),
			"error":     wrapped.Error(),
		})
		return wrapped
	}
	return o.runFromState(ctx, projectID, fromState)
}

func (o *Orchestrator) runFromState(ctx context.Context, projectID string, fromState State) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	fromIdx, ok := stateIndex(fromState)
	if !ok {
		err := fmt.Errorf("pipeline.Orchestrator: unknown state: %q", fromState)
		o.emit("pipeline:error", map[string]interface{}{
			"projectID": projectID,
			"step":      string(fromState),
			"error":     err.Error(),
		})
		return err
	}

	if fromIdx == len(StateOrder)-1 {
		o.emit("pipeline:complete", map[string]interface{}{"projectID": projectID})
		return nil
	}

	for i := fromIdx; i < len(StateOrder)-1; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		stepState := StateOrder[i]
		stepFn := o.steps[stepState]
		if stepFn == nil {
			err := fmt.Errorf("pipeline.Orchestrator: no step registered for state %q", stepState)
			o.emit("pipeline:error", map[string]interface{}{
				"projectID": projectID,
				"step":      string(stepState),
				"error":     err.Error(),
			})
			return err
		}

		o.emit("pipeline:step:start", map[string]interface{}{
			"projectID": projectID,
			"step":      string(stepState),
		})

		next, err := o.executeStep(ctx, projectID, stepFn)
		if err != nil {
			wrapped := fmt.Errorf("pipeline.Orchestrator: step %s failed: %w", stepState, err)
			o.emit("pipeline:error", map[string]interface{}{
				"projectID": projectID,
				"step":      string(stepState),
				"error":     wrapped.Error(),
			})
			return wrapped
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !IsValidTransition(stepState, next) {
			err := fmt.Errorf("pipeline.Orchestrator: invalid transition %q -> %q", stepState, next)
			o.emit("pipeline:error", map[string]interface{}{
				"projectID": projectID,
				"step":      string(stepState),
				"error":     err.Error(),
			})
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := o.config.StateStore.SetState(projectID, next); err != nil {
			wrapped := fmt.Errorf("pipeline.Orchestrator: persist state after %s (%s): %w", stepState, next, err)
			o.emit("pipeline:error", map[string]interface{}{
				"projectID": projectID,
				"step":      string(stepState),
				"error":     wrapped.Error(),
			})
			return wrapped
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
		o.emit("pipeline:step:complete", map[string]interface{}{
			"projectID": projectID,
			"step":      string(stepState),
		})

		if idx, ok := stateIndex(next); ok && idx != i+1 {
			i = idx - 1
		}
	}

	o.emit("pipeline:complete", map[string]interface{}{"projectID": projectID})
	return nil
}

func (o *Orchestrator) emit(event string, data interface{}) {
	if o == nil {
		return
	}
	if o.config.Emitter == nil {
		return
	}
	o.config.Emitter.Emit(event, data)
}

func stateIndex(s State) (int, bool) {
	for i, st := range StateOrder {
		if st == s {
			return i, true
		}
	}
	return -1, false
}

func (o *Orchestrator) executeStep(ctx context.Context, projectID string, fn StepFunc) (State, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	next, err := fn(ctx, projectID)
	if err == nil {
		return next, nil
	}
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if !isRetryableForOrchestrator(err) {
		return "", err
	}

	effective := o.config.RetryConfig.withDefaults()
	if effective.MaxAttempts <= 1 {
		return "", err
	}
	retryCfg := effective
	retryCfg.MaxAttempts = effective.MaxAttempts - 1

	return Retry[State](ctx, retryCfg, func() (State, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return fn(ctx, projectID)
	})
}

func isRetryableForOrchestrator(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if IsRetryable(e) {
			return true
		}
	}
	return false
}
