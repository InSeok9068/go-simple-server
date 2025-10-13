package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	hostinstrumentation "go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	runtimeinstrumentation "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var shutdown func(context.Context) error

func InitTracer() {
	ctx := context.Background()
	shutdown = nil
	defer setupHTTPInstrumentation()

	serviceName := os.Getenv("SERVICE_NAME")
	deployEnv := os.Getenv("ENV")

	resourceAttributes := []attribute.KeyValue{
		semconv.ServiceName(serviceName),
	}
	if deployEnv != "" {
		resourceAttributes = append(resourceAttributes, semconv.DeploymentEnvironment(deployEnv))
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(resourceAttributes...),
	)
	if err != nil {
		slog.Error("리소스 초기화 실패", "error", err)
		return
	}

	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if licenseKey == "" {
		initLocalProviders(res)
		slog.Warn("NewRelic 라이선스를 찾지 못해 로컬 OTEL 프로바이더만 활성화합니다")
		return
	}

	localProvider := os.Getenv("LOCAL_PROVIDER")
	if localProvider == "true" {
		initLocalProviders(res)
		return
	}

	endpoint := os.Getenv("NEW_RELIC_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://otlp.nr-data.net:4318"
	}

	shutdowns, ok := setupRemoteProviders(ctx, res, endpoint, licenseKey)
	if !ok {
		return
	}

	if len(shutdowns) == 0 {
		return
	}

	shutdown = func(ctx context.Context) error {
		var firstErr error
		for i := len(shutdowns) - 1; i >= 0; i-- {
			if err := shutdowns[i](ctx); err != nil {
				if firstErr == nil {
					firstErr = err
				} else {
					slog.Warn("추가 종료 오류 발생", "error", err)
				}
			}
		}
		return firstErr
	}

	slog.Info("OpenTelemetry 추적/메트릭을 NewRelic으로 전송합니다", "endpoint", endpoint, "service", serviceName, "env", deployEnv)
}

func setupRemoteProviders(
	ctx context.Context,
	res *resource.Resource,
	endpoint, licenseKey string,
) ([]func(context.Context) error, bool) {
	var shutdowns []func(context.Context) error
	registerShutdown := func(fn func(context.Context) error) {
		shutdowns = append(shutdowns, fn)
	}

	tp, traceErr := initTraceProvider(ctx, res, endpoint, licenseKey, registerShutdown)
	mp, metricErr := initMetricProvider(ctx, res, endpoint, licenseKey, registerShutdown)

	if traceErr != nil && metricErr != nil {
		slog.Error("OTel 원격 초기화 실패, 로컬 프로바이더로 대체합니다", "trace_error", traceErr, "metric_error", metricErr)
		initLocalProviders(res)
		return nil, false
	}
	if traceErr != nil {
		slog.Error("추적 프로바이더 초기화 실패, 로컬 트레이서로 대체합니다", "error", traceErr)
		initLocalTraceProvider(res, registerShutdown)
	}
	if metricErr != nil {
		slog.Error("메트릭 프로바이더 초기화 실패, 로컬 메트릭으로 대체합니다", "error", metricErr)
		initLocalMetricProvider(res, registerShutdown)
	}
	// 둘 다 성공한 경우만 원격 경로 검증 및 런타임 수집 시작
	if traceErr == nil && metricErr == nil {
		if ok := ensureRemoteReady(ctx, tp, mp, shutdowns, res); !ok {
			return nil, false
		}
	}

	return shutdowns, true
}

func initTraceProvider(
	ctx context.Context,
	res *resource.Resource,
	endpoint, licenseKey string,
	registerShutdown func(func(context.Context) error),
) (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpointURL(endpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"api-key": licenseKey,
		}),
	)

	exporterCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exporter, err := otlptrace.New(exporterCtx, client)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTracerProvider(tp)
	registerShutdown(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	})

	return tp, nil
}

func initMetricProvider(
	ctx context.Context,
	res *resource.Resource,
	endpoint, licenseKey string,
	registerShutdown func(func(context.Context) error),
) (*metric.MeterProvider, error) {
	exporterCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exporter, err := otlpmetrichttp.New(
		exporterCtx,
		otlpmetrichttp.WithEndpointURL(endpoint),
		otlpmetrichttp.WithHeaders(map[string]string{
			"api-key": licenseKey,
		}),
	)
	if err != nil {
		return nil, err
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(15*time.Second))

	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)

	otel.SetMeterProvider(mp)
	registerShutdown(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return mp.Shutdown(ctx)
	})

	return mp, nil
}

func initLocalTraceProvider(res *resource.Resource, registerShutdown func(func(context.Context) error)) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	if registerShutdown != nil {
		registerShutdown(func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return tp.Shutdown(ctx)
		})
	}

	return tp
}

func initLocalMetricProvider(res *resource.Resource, registerShutdown func(func(context.Context) error)) *metric.MeterProvider {
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	if registerShutdown != nil {
		registerShutdown(func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return mp.Shutdown(ctx)
		})
	}

	return mp
}

func initLocalProviders(res *resource.Resource) {
	tp := initLocalTraceProvider(res, nil)
	mp := initLocalMetricProvider(res, nil)

	shutdown = func(ctx context.Context) error {
		var firstErr error

		if tp != nil {
			if err := func() error {
				ctxTrace, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				return tp.Shutdown(ctxTrace)
			}(); err != nil {
				firstErr = err
			}
		}

		if mp != nil {
			if err := func() error {
				ctxMetric, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				return mp.Shutdown(ctxMetric)
			}(); err != nil {
				if firstErr == nil {
					firstErr = err
				} else {
					slog.Warn("메트릭 종료 실패", "error", err)
				}
			}
		}

		return firstErr
	}
}

func setupHTTPInstrumentation() {
	base := http.DefaultTransport
	if base == nil {
		base = &http.Transport{}
	}
	if _, ok := base.(*otelhttp.Transport); ok {
		http.DefaultTransport = base
		return
	}
	http.DefaultTransport = otelhttp.NewTransport(base)
}

func ShutdownTracer(ctx context.Context) {
	if shutdown != nil {
		if err := shutdown(ctx); err != nil {
			slog.Error("OTel 리소스 종료 실패", "error", err)
		}
	}
}

// ensureRemoteReady는 핑 전송으로 원격 경로를 검증하고,
// 성공 시 런타임/호스트 메트릭 수집을 시작합니다. 실패 시 로컬 폴백합니다.
func ensureRemoteReady(
	ctx context.Context,
	tp *sdktrace.TracerProvider,
	mp *metric.MeterProvider,
	shutdowns []func(context.Context) error,
	res *resource.Resource,
) bool {
	if err := pingOTel(ctx, tp, mp); err != nil {
		slog.Error("OTel 원격 전송 핑 실패, 로컬로 폴백합니다", "error", err)
		for i := len(shutdowns) - 1; i >= 0; i-- {
			_ = shutdowns[i](context.Background())
		}
		initLocalProviders(res)
		return false
	}

	meterProvider := otel.GetMeterProvider()
	if err := runtimeinstrumentation.Start(
		runtimeinstrumentation.WithMeterProvider(meterProvider),
		runtimeinstrumentation.WithMinimumReadMemStatsInterval(15*time.Second),
	); err != nil {
		slog.Warn("런타임 메트릭 수집 초기화 실패", "error", err)
	}
	if err := hostinstrumentation.Start(
		hostinstrumentation.WithMeterProvider(meterProvider),
	); err != nil {
		slog.Warn("호스트 메트릭 수집 초기화 실패", "error", err)
	}
	return true
}

// pingOTel은 초기화 직후 원격 수집기로의 전송 가능 여부를 확인하기 위해
// 짧은 스팬/메트릭을 전송하고 즉시 Flush 합니다. 실패 시 에러를 반환합니다.
func pingOTel(ctx context.Context, tp *sdktrace.TracerProvider, mp *metric.MeterProvider) error {
	// Trace ping
	{
		tctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		tracer := tp.Tracer("otel-init")
		_, span := tracer.Start(tctx, "otel.init.ping")
		span.SetAttributes(attribute.String("ping", "true"))
		span.End()
		if err := tp.ForceFlush(tctx); err != nil {
			return fmt.Errorf("trace flush 실패: %w", err)
		}
	}

	// Metric ping
	{
		mctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		meter := mp.Meter("otel-init")
		if counter, err := meter.Int64Counter("otel_init_ping"); err == nil {
			counter.Add(mctx, 1)
		}
		if err := mp.ForceFlush(mctx); err != nil {
			return fmt.Errorf("metric flush 실패: %w", err)
		}
	}

	slog.Debug("OTel 원격 핑 성공")
	return nil
}
