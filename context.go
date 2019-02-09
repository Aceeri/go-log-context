package logContext

import (
	"context"
	"errors"
	"time"
)

var Sleep = errors.New("sleep")
var Fatal = errors.New("fatal")

// Context is the state/situation of where the code is being called from.
// Contains various structures for hierarchical logging, tracking, and metrics.
type Context struct {
	Name      string
	ManagerId *int64 // Manager id

	Context  context.Context
	Logger   Logger
	Tracker  Tracker
	Triggers []Trigger
	Metrics  Metrics
}

// Fork creates a second context and forks the underlying structures, such as logging.
func (context *Context) Fork(name string) Context {
	copy := *context
	copy.Context = context
	copy.Name = name
	copy.Logger = copy.Logger.Fork(name)
	return copy
}

// ChildTimeout wraps the context with a timeout context.
func (ctx *Context) ChildTimeout(timeout time.Duration) (Context, context.CancelFunc) {
	copy := *ctx
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	copy.Context = timeoutCtx
	return copy, cancel
}

// ChildTimeout wraps the context with a cancelling context.
func (ctx *Context) ChildCancel() (Context, context.CancelFunc) {
	copy := *ctx
	cancelCtx, cancel := context.WithCancel(ctx)
	copy.Context = cancelCtx
	return copy, cancel
}

// SetTracker updates the associated a tracker to the context at a location.
func (context *Context) SetTracker(location ...string) {
	context.Tracker.Set(location...)
}

func (context *Context) GetTracker() Tracker {
	return context.Tracker
}

func (context *Context) SetDebug(debug bool) {
	context.Logger.SetDebug(debug)
}

func (context *Context) GetDebug() bool {
	return context.Logger.GetDebug()
}

func (context *Context) Dlog(format string, args ...interface{}) {
	context.Logger.Dlog(format, args...)
}

// Log forwards to the inner logger.
func (context *Context) Log(format string, args ...interface{}) {
	context.Logger.Log(format, args...)
}

func (context *Context) RawLog(format string, args ...interface{}) {
	context.Logger.RawLog(format, args...)
}

// Elog forwards to the inner logger error logging.
func (context *Context) Elog(format string, args ...interface{}) {
	context.Logger.Elog(format, args...)
}

func (context *Context) Trigger() {
	for _, event := range context.Triggers {
		trigger(event)
	}
}

func (context *Context) GetLogger() Logger {
	return context.Logger
}

// Implement golang `context` interface
func (ctx Context) Deadline() (time.Time, bool) {
	if ctx.Context == nil {
		ctx.Context = context.Background()
	}

	return ctx.Context.Deadline()
}

func (ctx Context) Done() <-chan struct{} {
	if ctx.Context == nil {
		ctx.Context = context.Background()
	}

	return ctx.Context.Done()
}

func (ctx Context) Err() error {
	if ctx.Context == nil {
		ctx.Context = context.Background()
	}

	return ctx.Context.Err()
}

func (ctx Context) Value(key interface{}) interface{} {
	if ctx.Context == nil {
		ctx.Context = context.Background()
	}

	return ctx.Context.Value(key)
}

type Trigger = chan struct{}

func trigger(c Trigger) {
	select {
	case c <- struct{}{}:
	default:
	}
}
