package zerolog

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

const here = "github.com/prixfixeco/api_server/"

func init() {
	zerolog.CallerSkipFrameCount += 2
	zerolog.DisableSampling(true)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		return strings.TrimPrefix(file, here) + ", line " + strconv.Itoa(line)
	}
}

// logger is our log wrapper.
type zerologLogger struct {
	requestIDFunc logging.RequestIDFunc
	logger        zerolog.Logger
}

// buildZerologger builds a new zerologger.
func buildZerologger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.InfoLevel)
}

// NewZerologLogger builds a new zerologLogger.
func NewZerologLogger() logging.Logger {
	return &zerologLogger{logger: buildZerologger()}
}

// WithName is our obligatory contract fulfillment function.
// Zerolog doesn't support named loggers :( so we have this workaround.
func (l *zerologLogger) WithName(name string) logging.Logger {
	l2 := l.logger.With().Str(logging.LoggerNameKey, name).Logger()
	return &zerologLogger{logger: l2}
}

// SetLevel sets the log level for our zerolog logger.
func (l *zerologLogger) SetLevel(level logging.Level) {
	var lvl zerolog.Level

	switch level {
	case logging.InfoLevel:
		lvl = zerolog.InfoLevel
	case logging.DebugLevel:
		l.logger = l.logger.With().Logger()
		lvl = zerolog.DebugLevel
	case logging.WarnLevel:
		l.logger = l.logger.With().Caller().Logger()
		lvl = zerolog.WarnLevel
	case logging.ErrorLevel:
		l.logger = l.logger.With().Caller().Logger()
		lvl = zerolog.ErrorLevel
	default:
		lvl = zerolog.InfoLevel
	}

	l.logger = l.logger.Level(lvl)
}

// SetRequestIDFunc sets the request ID retrieval function.
func (l *zerologLogger) SetRequestIDFunc(f logging.RequestIDFunc) {
	if f != nil {
		l.requestIDFunc = f
	}
}

// Info satisfies our contract for the logging.Logger Info method.
func (l *zerologLogger) Info(input string) {
	l.logger.Info().Msg(input)
}

// Debug satisfies our contract for the logging.Logger Debug method.
func (l *zerologLogger) Debug(input string) {
	l.logger.Debug().Msg(input)
}

// Error satisfies our contract for the logging.Logger Error method.
func (l *zerologLogger) Error(err error, whatWasHappeningWhenErrorOccurred string) {
	if err != nil {
		l.logger.Error().Stack().Caller().Msg(fmt.Sprintf("error %s: %s", whatWasHappeningWhenErrorOccurred, err.Error()))
	}
}

// Printf satisfies our contract for the logging.Logger Printf method.
func (l *zerologLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

// Clone satisfies our contract for the logging.Logger WithValue method.
func (l *zerologLogger) Clone() logging.Logger {
	l2 := l.logger.With().Logger()
	return &zerologLogger{logger: l2}
}

// WithValue satisfies our contract for the logging.Logger WithValue method.
func (l *zerologLogger) WithValue(key string, value interface{}) logging.Logger {
	l2 := l.logger.With().Interface(key, value).Logger()
	return &zerologLogger{logger: l2}
}

// WithValues satisfies our contract for the logging.Logger WithValues method.
func (l *zerologLogger) WithValues(values map[string]interface{}) logging.Logger {
	var l2 = l.logger.With().Logger()

	for key, val := range values {
		l2 = l2.With().Interface(key, val).Logger()
	}

	return &zerologLogger{logger: l2}
}

// WithError satisfies our contract for the logging.Logger WithError method.
func (l *zerologLogger) WithError(err error) logging.Logger {
	l2 := l.logger.With().Err(err).Logger()
	return &zerologLogger{logger: l2}
}

// WithSpan satisfies our contract for the logging.Logger WithSpan method.
func (l *zerologLogger) WithSpan(span trace.Span) logging.Logger {
	spanCtx := span.SpanContext()
	spanID := spanCtx.SpanID().String()
	traceID := spanCtx.TraceID().String()

	l2 := l.logger.With().Str(keys.SpanIDKey, spanID).Str(keys.TraceIDKey, traceID).Logger()

	return &zerologLogger{logger: l2}
}

func (l *zerologLogger) attachRequestToLog(req *http.Request) zerolog.Logger {
	if req != nil {
		l2 := l.logger.With().
			Str("method", req.Method).
			Logger()

		if req.URL != nil {
			l2 = l2.With().Str("path", req.URL.Path).Logger()
			if req.URL.RawQuery != "" {
				l2 = l2.With().Str(keys.URLQueryKey, req.URL.RawQuery).Logger()
			}
		}

		if l.requestIDFunc != nil {
			if reqID := l.requestIDFunc(req); reqID != "" {
				l2 = l2.With().Str("request.id", reqID).Logger()
			}
		}

		return l2
	}

	return l.logger
}

// WithRequest satisfies our contract for the logging.Logger WithRequest method.
func (l *zerologLogger) WithRequest(req *http.Request) logging.Logger {
	return &zerologLogger{logger: l.attachRequestToLog(req)}
}

// WithResponse satisfies our contract for the logging.Logger WithResponse method.
func (l *zerologLogger) WithResponse(res *http.Response) logging.Logger {
	l2 := l.logger.With().Logger()
	if res != nil {
		l2 = l.attachRequestToLog(res.Request).With().Int(keys.ResponseStatusKey, res.StatusCode).Logger()
	}

	return &zerologLogger{logger: l2}
}
