package debug

import (
	"expvar"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	reqCount       = expvar.NewInt("req_count")
	errCount       = expvar.NewInt("err_count")
	totalLatencyMS = expvar.NewInt("total_latency_ms")
)

// MetricsMiddleware는 요청 수, 에러 수, 응답 시간을 기록한다.
func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		dur := time.Since(start).Milliseconds()
		reqCount.Add(1)
		totalLatencyMS.Add(dur)
		if err != nil || c.Response().Status >= 400 {
			errCount.Add(1)
		}
		return err
	}
}
