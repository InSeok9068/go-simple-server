// SPDX-License-Identifier: Apache-2.0
// 단일 파일: Echo용 OTEL 미들웨어 (원본 구조 최대 유지, c.Error 제거, internal/semconv 대체)
// https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/labstack/echo/otelecho
package middleware // 사용처에서 import "your/module/path/otelecho"

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// ─────────────────────────────────────────────────────────────────────────────
// 공개 API: 옵션/config/버전 (원본과 동일 패턴, 명칭만 otelechoConfig)
// ─────────────────────────────────────────────────────────────────────────────

type Option interface{ apply(*otelechoConfig) }

type optFn func(*otelechoConfig)

func (f optFn) apply(c *otelechoConfig) { f(c) }

// otelechoConfig: 원본 config에 맞춰 필드 구성
type otelechoConfig struct {
	TracerProvider        trace.TracerProvider
	Propagators           propagation.TextMapPropagator
	MeterProvider         metric.MeterProvider
	Skipper               middleware.Skipper
	SpanNameFormatter     func(echo.Context) string
	MetricAttributeFn     func(*http.Request) []attribute.KeyValue
	EchoMetricAttributeFn func(echo.Context) []attribute.KeyValue
}

// 옵션들
func WithTracerProvider(tp trace.TracerProvider) Option {
	return optFn(func(c *otelechoConfig) { c.TracerProvider = tp })
}
func WithPropagators(p propagation.TextMapPropagator) Option {
	return optFn(func(c *otelechoConfig) { c.Propagators = p })
}
func WithMeterProvider(mp metric.MeterProvider) Option {
	return optFn(func(c *otelechoConfig) { c.MeterProvider = mp })
}
func WithSkipper(s middleware.Skipper) Option {
	return optFn(func(c *otelechoConfig) { c.Skipper = s })
}
func WithSpanNameFormatter(fn func(echo.Context) string) Option {
	return optFn(func(c *otelechoConfig) { c.SpanNameFormatter = fn })
}
func WithMetricAttributeFn(fn func(*http.Request) []attribute.KeyValue) Option {
	return optFn(func(c *otelechoConfig) { c.MetricAttributeFn = fn })
}
func WithEchoMetricAttributeFn(fn func(echo.Context) []attribute.KeyValue) Option {
	return optFn(func(c *otelechoConfig) { c.EchoMetricAttributeFn = fn })
}

// Version: 원본은 빌드 시 주입되지만, 여기선 상수로 처리
func Version() string { return "v0-local" }

// ─────────────────────────────────────────────────────────────────────────────
// 내부 semconv 대체 (원본 internal/semconv를 단일 파일 안에 최소 재현)
// ─────────────────────────────────────────────────────────────────────────────

type httpServerSemconv struct {
	meter metric.Meter
}

func newHTTPServerSemconv(m metric.Meter) *httpServerSemconv {
	return &httpServerSemconv{meter: m}
}

// 요청 시 붙일 속성들
type requestTraceAttrsOpts struct{}

func (s *httpServerSemconv) RequestTraceAttrs(
	service string,
	r *http.Request,
	_ requestTraceAttrsOpts,
) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(service),
		semconv.NetworkProtocolNameKey.String("http"),
		semconv.HTTPRequestMethodKey.String(r.Method),
		semconv.URLFullKey.String(r.URL.String()),
		semconv.UserAgentOriginalKey.String(r.UserAgent()),
	}
	if r.URL != nil {
		attrs = append(attrs, attribute.String("http.target", r.URL.Path))
	}
	if host := r.Host; host != "" {
		attrs = append(attrs, semconv.ServerAddressKey.String(host))
	}
	return attrs
}

func (s *httpServerSemconv) Route(path string) attribute.KeyValue {
	return semconv.HTTPRouteKey.String(path)
}

func (s *httpServerSemconv) Status(status int) (codes.Code, string) {
	// 원본 정책: 5xx만 Error 처리(4xx는 Unset)
	if status >= 500 {
		return codes.Error, http.StatusText(status)
	}
	return codes.Unset, http.StatusText(status)
}

type responseTelemetry struct {
	StatusCode int
	WriteBytes int64
}

func (s *httpServerSemconv) ResponseTraceAttrs(rt responseTelemetry) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.HTTPResponseStatusCodeKey.Int(rt.StatusCode),
		attribute.Int64("http.response.body.size", rt.WriteBytes),
	}
}

type metricAttributes struct {
	Req                  *http.Request
	StatusCode           int
	AdditionalAttributes []attribute.KeyValue
}
type metricData struct {
	RequestSize int64
	ElapsedTime float64 // ms
}
type serverMetricData struct {
	ServerName       string
	ResponseSize     int64
	MetricAttributes metricAttributes
	MetricData       metricData
}

// 최소 구현: 메트릭 기록은 no-op (필요 시 확장)
func (s *httpServerSemconv) RecordMetrics(ctx context.Context, _ serverMetricData) {
	// no-op
}

// ─────────────────────────────────────────────────────────────────────────────
// 미들웨어 본체 (원본과 동일 흐름, 단 err 처리에서 c.Error 제거)
// ─────────────────────────────────────────────────────────────────────────────

const (
	tracerKey = "otel-go-contrib-tracer-labstack-echo"
	// ScopeName는 원본과 동일하게 둠 (관측툴에서 식별 용이)
	ScopeName = "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// Middleware returns echo middleware which will trace incoming requests.
func OtelEchoMiddleware(service string, opts ...Option) echo.MiddlewareFunc {
	cfg := otelechoConfig{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		ScopeName,
		trace.WithInstrumentationVersion(Version()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	if cfg.MeterProvider == nil {
		cfg.MeterProvider = otel.GetMeterProvider()
	}
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}
	if cfg.SpanNameFormatter == nil {
		cfg.SpanNameFormatter = spanNameFormatter
	}

	meter := cfg.MeterProvider.Meter(
		ScopeName,
		metric.WithInstrumentationVersion(Version()),
	)
	semconvSrv := newHTTPServerSemconv(meter)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestStartTime := time.Now()
			if cfg.Skipper(c) {
				return next(c)
			}

			c.Set(tracerKey, tracer)

			request := c.Request()
			savedCtx := request.Context()
			defer func() {
				request = request.WithContext(savedCtx)
				c.SetRequest(request)
			}()

			// 상위 컨텍스트 추출
			ctx := cfg.Propagators.Extract(savedCtx, propagation.HeaderCarrier(request.Header))

			// 스팬 옵션 구성
			startOpts := []trace.SpanStartOption{
				trace.WithAttributes(
					semconvSrv.RequestTraceAttrs(service, request, requestTraceAttrsOpts{})...,
				),
				trace.WithSpanKind(trace.SpanKindServer),
			}
			if path := c.Path(); path != "" {
				rAttr := semconvSrv.Route(path)
				startOpts = append(startOpts, trace.WithAttributes(rAttr))
			}
			spanName := cfg.SpanNameFormatter(c)

			// 스팬 시작
			ctx, span := tracer.Start(ctx, spanName, startOpts...)
			defer span.End()

			// 컨텍스트 주입
			c.SetRequest(request.WithContext(ctx))

			// 실제 처리
			err := next(c)

			// ★ 변경점: 전역 에러 핸들러를 강제 호출하지 않음 (c.Error 제거)
			// 대신 스팬에만 에러 기록
			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
				span.RecordError(err)
			}

			// 상태/응답 속성
			status := c.Response().Status
			span.SetStatus(semconvSrv.Status(status))
			span.SetAttributes(semconvSrv.ResponseTraceAttrs(responseTelemetry{
				StatusCode: status,
				WriteBytes: c.Response().Size,
			})...)

			// 추가 속성 훅
			var additionalAttributes []attribute.KeyValue
			if path := c.Path(); path != "" {
				additionalAttributes = append(additionalAttributes, semconvSrv.Route(path))
			}
			if cfg.MetricAttributeFn != nil {
				additionalAttributes = append(additionalAttributes, cfg.MetricAttributeFn(request)...)
			}
			if cfg.EchoMetricAttributeFn != nil {
				additionalAttributes = append(additionalAttributes, cfg.EchoMetricAttributeFn(c)...)
			}

			// 메트릭 기록 (현재 no-op, 인터페이스 유지)
			semconvSrv.RecordMetrics(ctx, serverMetricData{
				ServerName:   service,
				ResponseSize: c.Response().Size,
				MetricAttributes: metricAttributes{
					Req:                  request,
					StatusCode:           status,
					AdditionalAttributes: additionalAttributes,
				},
				MetricData: metricData{
					RequestSize: request.ContentLength,
					ElapsedTime: float64(time.Since(requestStartTime)) / float64(time.Millisecond),
				},
			})

			return err
		}
	}
}

// 기본 스팬 이름 포맷터 (원본과 동일)
func spanNameFormatter(c echo.Context) string {
	method, path := strings.ToUpper(c.Request().Method), c.Path()
	if !slices.Contains([]string{
		http.MethodGet, http.MethodHead,
		http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions,
		http.MethodTrace,
	}, method) {
		method = "HTTP"
	}
	if path != "" {
		return method + " " + path
	}
	return method
}
