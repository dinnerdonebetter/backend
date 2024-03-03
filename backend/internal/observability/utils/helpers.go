package o11yutils

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"go.opentelemetry.io/otel/trace"
)

func ObserveValue[T any](logger logging.Logger, span trace.Span, key string, value T) (logging.Logger, trace.Span) {
	tracing.AttachToSpan(span, key, value)
	if logger != nil {
		logger = logger.WithValue(key, value)
	}

	return logger, span
}
