package handler

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
)

func init() {
	tracer = otel.GetTracerProvider().Tracer("internal/server/handler")
}
