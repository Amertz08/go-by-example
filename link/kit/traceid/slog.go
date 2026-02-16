package traceid

import (
	"context"
	"log/slog"
)

// LogKey is the key used to store the trace ID in the log record.
const LogKey = "trace_id"

// LogHandler adds trace IDs form [context.Context] to [slog.Record].
type LogHandler struct {
	slog.Handler // here is an example of field embedding
	LogKey       string
}

// NewLogHandler returns a new [LogHandler].
func NewLogHandler(h slog.Handler) slog.Handler {
	return &LogHandler{
		Handler: h,
		LogKey:  LogKey,
	}
}

// Handle may add a [slog.Attr] to the [slog.Record].
func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	// If the trace ID is present in the context, add it to the log record.
	if id, ok := FromContext(ctx); ok {
		r = r.Clone()
		r.AddAttrs(slog.String(h.LogKey, id))
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new [LogHandler] with the provided attributes.
func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLogHandler(h.Handler.WithAttrs(attrs))
}

// WithGroup returns a new [LogHandler] with the provided group.
func (h *LogHandler) WithGroup(group string) slog.Handler {
	return NewLogHandler(h.Handler.WithGroup(group))
}
