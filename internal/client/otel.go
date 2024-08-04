package client

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tp     = otel.GetTracerProvider()
	tracer trace.Tracer
)

func init() {
	tracer = otel.GetTracerProvider().Tracer("internal/client")
}
