package observer

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NOTE: I didn't develop `observer` specifically for this example, but it's something that I tend to init on the simple projects.

type Observer interface {
	ChildAction(ctx context.Context, tracer trace.Tracer, funcName ...string) (Observer, context.Context)
	Action(ctx context.Context, tracer trace.Tracer, funcName ...string) (Observer, context.Context)
	Str(key, val string) Observer
	UUID(key string, val uuid.UUID) Observer
	Int(key string, val int) Observer
	Bool(key string, val bool) Observer
	Len(key string, val any) Observer
	Any(key string, val any) Observer
	Int64(key string, val int64) Observer
	Float64(key string, val float64) Observer
	Close()
	Info(msg ...string)
	Err(err error) Observer
	Error(msg ...string)
	Fatal(msg ...string)
	Warn(msg ...string)
	Debug(msg ...string)
}

// observerImpl is a utility struct for logging and tracing. It encapsulates a Zerolog logger
// and an OpenTelemetry span for contextual logging and tracing.
type observerImpl struct {
	logger      zerolog.Logger
	span        trace.Span
	environment string
}

const (
	actionKey = "action"
)

var globalObserver Observer

// InitObserver initializes the global observer with the given service name, OpenTelemetry digester URL, and environment.
func InitObserver(serviceName string, otelDigesterURL string, environment string) {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	newLogger := zerolog.New(output).With().Timestamp().Logger()

	if otelDigesterURL != "" {
		InitTracer(serviceName, otelDigesterURL)
	}

	globalObserver = observerImpl{
		logger:      newLogger,
		environment: environment,
	}
}

// Action initializes an Observer with a custom tracer and a given action.
func Action(ctx context.Context, tracer trace.Tracer, funcName ...string) (Observer, context.Context) {
	return globalObserver.Action(ctx, tracer, funcName...)
}

// ChildAction creates a child span from the current span for a given action.
func (o observerImpl) ChildAction(ctx context.Context, tracer trace.Tracer, funcName ...string) (Observer, context.Context) {
	action := initActionFromOptionalInput(3, funcName...)

	childCtx, childSpan := tracer.Start(ctx, action)
	childObserver := observerImpl{
		logger:      o.logger,
		span:        childSpan,
		environment: o.environment,
	}

	sc := childSpan.SpanContext()
	if sc.IsValid() {
		traceID := sc.TraceID().String()
		childObserver.logger = o.logger.With().
			Str(actionKey, action).
			Str("trace_id", traceID).
			Str("env", o.environment).
			Logger()
	}

	return childObserver, childCtx
}

// Action creates or updates the Observer's span using the provided tracer, context, and action.
func (o observerImpl) Action(ctx context.Context, tracer trace.Tracer, funcName ...string) (Observer, context.Context) {
	action := initActionFromOptionalInput(4, funcName...)

	tCtx, span := tracer.Start(ctx, action)
	o.span = span

	sc := o.span.SpanContext()
	if !sc.IsValid() {
		o.logger = o.logger.With().Str(actionKey, action).Logger()
		return o, tCtx
	}

	traceID := sc.TraceID().String()

	o.logger = o.logger.With().
		Str(actionKey, action).
		Str("trace_id", traceID).
		Str("env", o.environment).
		Logger()

	return o, tCtx
}

func initActionFromOptionalInput(depth int, funcName ...string) string {
	var action string
	if len(funcName) == 0 || funcName[0] == "" {
		action = getCallingFuncName(depth)
	} else {
		action = funcName[0]
	}

	return action
}

// getCallingFuncName fetches the name of the calling function.
func getCallingFuncName(depth int) string {
	const unknown = "unknown caller"
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return unknown
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return unknown
	}

	funcName := fn.Name()
	if lastSlash := strings.LastIndex(funcName, "/"); lastSlash != -1 {
		funcName = funcName[lastSlash+1:]
	}
	if lastDot := strings.LastIndex(funcName, "."); lastDot != -1 {
		funcName = funcName[lastDot+1:]
	}

	return funcName
}

// Str adds a string attribute to both the logger and the span.
func (o observerImpl) Str(key, val string) Observer {
	o.logger = o.logger.With().Str(key, val).Logger()
	o.span.SetAttributes(attribute.String(key, val))

	return o
}

// UUID adds a UUID attribute to both the logger and the span.
func (o observerImpl) UUID(key string, val uuid.UUID) Observer {
	o.logger = o.logger.With().Str(key, val.String()).Logger()
	o.span.SetAttributes(attribute.String(key, val.String()))

	return o
}

// Int adds an integer attribute to both the logger and the span.
func (o observerImpl) Int(key string, val int) Observer {
	o.logger = o.logger.With().Int(key, val).Logger()
	o.span.SetAttributes(attribute.Int(key, val))

	return o
}

// Bool adds a boolean attribute to both the logger and the span.
func (o observerImpl) Bool(key string, val bool) Observer {
	o.logger = o.logger.With().Bool(key, val).Logger()
	o.span.SetAttributes(attribute.Bool(key, val))

	return o
}

// Len adds the length of a collection as an attribute to both the logger and the span.
func (o observerImpl) Len(key string, val any) Observer {
	var length int
	switch v := val.(type) {
	case string:
		length = len(v)
	case []any:
		length = len(v)
	case map[any]any:
		length = len(v)
	case chan any:
		length = len(v)
	default:
		o.logger.Warn().Msg("Attempted to take length of a type that does not support it")
		return o
	}
	o.logger = o.logger.With().Int(key, length).Logger()
	o.span.SetAttributes(attribute.Int(key, length))

	return o
}

// Any adds a generic attribute to both the logger and the span.
func (o observerImpl) Any(key string, val any) Observer {
	o.logger = o.logger.With().Interface(key, val).Logger()

	switch v := val.(type) {
	case string:
		o.span.SetAttributes(attribute.String(key, v))
	case int:
		o.span.SetAttributes(attribute.Int(key, v))
	case int64:
		o.span.SetAttributes(attribute.Int64(key, v))
	case bool:
		o.span.SetAttributes(attribute.Bool(key, v))
	case float64:
		o.span.SetAttributes(attribute.Float64(key, v))
	case []string:
		o.span.SetAttributes(attribute.StringSlice(key, v))
	case uuid.UUID:
		o.span.SetAttributes(attribute.String(key, v.String()))
	default:
		o.span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", val)))
	}

	return o
}

// Int64 adds a 64-bit integer attribute to both the logger and the span.
func (o observerImpl) Int64(key string, val int64) Observer {
	o.logger = o.logger.With().Int64(key, val).Logger()
	o.span.SetAttributes(attribute.Int64(key, val))

	return o
}

// Float64 adds a 64-bit floating-point attribute to both the logger and the span.
func (o observerImpl) Float64(key string, val float64) Observer {
	o.logger = o.logger.With().Float64(key, val).Logger()
	o.span.SetAttributes(attribute.Float64(key, val))

	return o
}

// Close ends the current span. It MUST be called when the operation represented by the span is completed.
func (o observerImpl) Close() {
	o.span.End()
}

// Info logs an informational message using the observer's logger.
func (o observerImpl) Info(msg ...string) {
	if len(msg) == 0 {
		o.logger.Info().Msg("")
		return
	}
	o.logger.Info().Msg(strings.Join(msg, " "))
}

// Err adds an error as a field to the logger and records it in the span.
func (o observerImpl) Err(err error) Observer {
	o.logger = o.logger.With().Err(err).Logger()
	o.span.RecordError(err)

	return o
}

// Error logs an error message using the observer's logger.
func (o observerImpl) Error(msg ...string) {
	if len(msg) == 0 {
		o.logger.Error().Msg("")
		return
	}
	o.logger.Error().Msg(strings.Join(msg, " "))
}

func (o observerImpl) Fatal(msg ...string) {
	if len(msg) == 0 {
		o.logger.Fatal().Msg("")
		return
	}
	o.logger.Fatal().Msg(strings.Join(msg, " "))
}

// Warn logs a warning message using the observer's logger.
func (o observerImpl) Warn(msg ...string) {
	if len(msg) == 0 {
		o.logger.Warn().Msg("")
		return
	}
	o.logger.Warn().Msg(strings.Join(msg, " "))
}

// Debug logs a debug message using the observer's logger.
func (o observerImpl) Debug(msg ...string) {
	if len(msg) == 0 {
		o.logger.Debug().Msg("")
		return
	}
	o.logger.Debug().Msg(strings.Join(msg, " "))
}
