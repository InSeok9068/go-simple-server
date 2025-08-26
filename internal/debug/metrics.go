package debug

import (
	"expvar"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	reqCount       = expvar.NewInt("req_count")
	errCount       = expvar.NewInt("err_count")
	errCountByPath = expvar.NewMap("err_count_by_path")
	totalLatencyMS = expvar.NewInt("total_latency_ms")

	mu sync.Mutex

	// 최근 지연시간 샘플 (슬라이딩 버퍼)
	latencySamples [256]int64
	latencyPos     int
	latencyFilled  bool

	// 최근 1분 요청/에러 수
	secReqs  [60]int64
	secErrs  [60]int64
	secTimes [60]int64
)

// MetricsMiddleware는 요청 수, 에러 수, 응답 시간을 기록한다.
func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		dur := time.Since(start).Milliseconds()

		reqCount.Add(1)
		totalLatencyMS.Add(dur)

		isErr := err != nil || c.Response().Status >= 400
		if isErr {
			errCount.Add(1)
			errCountByPath.Add(c.Path(), 1)
		}

		recordMetrics(dur, isErr, time.Now())

		return err
	}
}

func recordMetrics(dur int64, isErr bool, t time.Time) {
	mu.Lock()
	// 지연시간 샘플 저장
	latencySamples[latencyPos] = dur
	latencyPos++
	if latencyPos >= len(latencySamples) {
		latencyPos = 0
		latencyFilled = true
	}

	// 최근 1분 요청/에러 수 기록
	sec := t.Unix()
	idx := int(sec % int64(len(secReqs)))
	if secTimes[idx] != sec {
		secTimes[idx] = sec
		secReqs[idx] = 0
		secErrs[idx] = 0
	}
	secReqs[idx]++
	if isErr {
		secErrs[idx]++
	}
	mu.Unlock()
}

func latencyPercentiles() (float64, float64) {
	mu.Lock()
	n := latencyPos
	if latencyFilled {
		n = len(latencySamples)
	}
	data := make([]int64, n)
	copy(data, latencySamples[:n])
	mu.Unlock()

	if n == 0 {
		return 0, 0
	}
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })
	p95 := float64(data[int(math.Ceil(0.95*float64(n))-1)])
	p99 := float64(data[int(math.Ceil(0.99*float64(n))-1)])
	return p95, p99
}

func oneMinuteCounts() (int64, int64) {
	mu.Lock()
	now := time.Now().Unix()
	var reqs, errs int64
	for i := 0; i < len(secReqs); i++ {
		if now-secTimes[i] < int64(len(secReqs)) {
			reqs += secReqs[i]
			errs += secErrs[i]
		}
	}
	mu.Unlock()
	return reqs, errs
}
