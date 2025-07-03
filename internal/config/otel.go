package config

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
)

var shutdown func(context.Context) error

func InitTracer() {
	ctx := context.Background()
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		slog.Error("트레이스 exporter 생성 실패", "error", err)
		return
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(os.Getenv("SERVICE_NAME")),
		),
	)
	if err != nil {
		slog.Error("리소스 생성 실패", "error", err)
		return
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	shutdown = tp.Shutdown
}

func ShutdownTracer(ctx context.Context) {
	if shutdown != nil {
		if err := shutdown(ctx); err != nil {
			slog.Error("트레이서 종료 실패", "error", err)
		}
	}
}
