package observability

import (
	"context"
	"time"
	"github.com/segmentio/ksuid"
)

// Trace represents a single trace
type Trace struct {
	ID        string
	ParentID  string
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Tags      map[string]string
	Children  []*Trace
}

// Tracer manages traces
type Tracer struct {
	logger *Logger
}

type traceKey struct{}

// NewTracer creates a new tracer
func NewTracer(logger *Logger) *Tracer {
	return &Tracer{
		logger: logger,
	}
}

// StartTrace starts a new trace and returns the context with trace
func (t *Tracer) StartTrace(ctx context.Context, name string) (context.Context, *Trace) {
	trace := &Trace{
		ID:        ksuid.New().String(),
		Name:      name,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
	}

	if parent, ok := ctx.Value(traceKey{}).(*Trace); ok {
		trace.ParentID = parent.ID
		parent.Children = append(parent.Children, trace)
	}

	return context.WithValue(ctx, traceKey{}, trace), trace
}

// EndTrace ends the trace and logs it
func (t *Tracer) EndTrace(ctx context.Context, trace *Trace) {
	trace.EndTime = time.Now()
	duration := trace.EndTime.Sub(trace.StartTime)

	// Log trace information
	t.logger.Info("Trace completed",
		LogFields{
			"trace_id":   trace.ID,
			"parent_id":  trace.ParentID,
			"name":       trace.Name,
			"duration":   duration.Milliseconds(),
			"tags":       trace.Tags,
			"start_time": trace.StartTime,
			"end_time":   trace.EndTime,
		})
}

// AddTag adds a tag to the current trace
func (t *Tracer) AddTag(ctx context.Context, key, value string) {
	if trace, ok := ctx.Value(traceKey{}).(*Trace); ok {
		trace.Tags[key] = value
	}
}

// GetCurrentTrace gets the current trace from context
func (t *Tracer) GetCurrentTrace(ctx context.Context) *Trace {
	if trace, ok := ctx.Value(traceKey{}).(*Trace); ok {
		return trace
	}
	return nil
}