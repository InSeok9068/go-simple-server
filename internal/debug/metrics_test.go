package debug

import (
	"testing"
	"time"
)

func resetMetrics() {
	mu.Lock()
	latencyPos = 0
	latencyFilled = false
	for i := range latencySamples {
		latencySamples[i] = 0
	}
	for i := range secReqs {
		secReqs[i] = 0
		secErrs[i] = 0
		secTimes[i] = 0
	}
	mu.Unlock()
	reqCount.Set(0)
	errCount.Set(0)
	totalLatencyMS.Set(0)
}

func TestLatencyPercentiles(t *testing.T) {
	resetMetrics()
	base := time.Now()
	for i := 1; i <= 100; i++ {
		recordMetrics(int64(i), false, base)
	}
	p95, p99 := latencyPercentiles()
	if p95 != 95 {
		t.Fatalf("p95=%v", p95)
	}
	if p99 != 99 {
		t.Fatalf("p99=%v", p99)
	}
}

func TestOneMinuteCounts(t *testing.T) {
	resetMetrics()
	now := time.Now()
	// 7 reqs 70s ago (ignored)
	old := now.Add(-70 * time.Second)
	for i := 0; i < 7; i++ {
		recordMetrics(0, false, old)
	}
	// 5 reqs 10s ago, 1 error
	past := now.Add(-10 * time.Second)
	for i := 0; i < 5; i++ {
		recordMetrics(0, i == 0, past)
	}
	// 10 reqs now, 2 errors
	for i := 0; i < 10; i++ {
		recordMetrics(0, i < 2, now)
	}
	reqs, errs := oneMinuteCounts()
	if reqs != 15 || errs != 3 {
		t.Fatalf("reqs=%d errs=%d", reqs, errs)
	}
}
