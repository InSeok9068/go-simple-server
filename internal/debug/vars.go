package debug

import (
	"context"
	"database/sql"
	"expvar"
	"fmt"
	"math"
	"net/http"
	"os"
	"runtime"
	"time"
)

type DBSnapshot struct {
	Open     int           `json:"open"`
	MaxOpen  int           `json:"max_open"`
	InUse    int           `json:"in_use"`
	Idle     int           `json:"idle"`
	WaitCnt  int64         `json:"wait_count"`
	WaitDur  time.Duration `json:"wait_duration"`
	WaitRate float64       `json:"wait_rate"`
	Pragma   any           `json:"pragma"`
}

type MemSnapshot struct {
	AllocMB      float64 `json:"alloc_mb"`
	TotalAllocMB float64 `json:"total_alloc_mb"`
	SysMB        float64 `json:"sys_mb"`
	HeapAllocMB  float64 `json:"heap_alloc_mb"`
	HeapSysMB    float64 `json:"heap_sys_mb"`
	NumGC        uint32  `json:"num_gc"`
	PauseTotalMS float64 `json:"pause_total_ms"`
	NextGCMb     float64 `json:"next_gc_mb"`
}

type RuntimeSnapshot struct {
	Goroutines int `json:"goroutines"`
}

type HTTPSnapshot struct {
	AvgResMS     float64          `json:"avg_res_ms"`
	P95ResMS     float64          `json:"p95_res_ms"`
	P99ResMS     float64          `json:"p99_res_ms"`
	RPS          float64          `json:"rps"`
	ErrorRate    float64          `json:"error_rate"`
	ErrRate1m    float64          `json:"err_rate_1m"`
	TotalReqs    int64            `json:"total_reqs"`
	TotalErrors  int64            `json:"total_errors"`
	ErrorsByPath map[string]int64 `json:"errors_by_path"`
}

type AppSnapshot struct {
	At      time.Time       `json:"at"`
	DB      DBSnapshot      `json:"db"`
	Mem     MemSnapshot     `json:"mem"`
	Runtime RuntimeSnapshot `json:"runtime"`
	HTTP    HTTPSnapshot    `json:"http"`
	Note    string          `json:"note"`
}

func Init(name string, db *sql.DB) {
	expvar.Publish(name, expvar.Func(func() any {
		return takeSnapshot(db) // 요청 들어올 때만 실행
	}))
}

func takeSnapshot(db *sql.DB) AppSnapshot {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	mem := MemSnapshot{
		AllocMB:      bytesToMB(ms.Alloc),
		TotalAllocMB: bytesToMB(ms.TotalAlloc),
		SysMB:        bytesToMB(ms.Sys),
		HeapAllocMB:  bytesToMB(ms.HeapAlloc),
		HeapSysMB:    bytesToMB(ms.HeapSys),
		NumGC:        ms.NumGC,
		PauseTotalMS: float64(ms.PauseTotalNs) / 1e6,
		NextGCMb:     bytesToMB(ms.NextGC),
	}

	rt := RuntimeSnapshot{
		Goroutines: runtime.NumGoroutine(),
	}

	s := db.Stats()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	pr := map[string]any{
		"journal_mode":       pragmaString(ctx, db, "journal_mode"),
		"synchronous":        pragmaEnum(ctx, db, "synchronous", map[int64]string{0: "OFF", 1: "NORMAL", 2: "FULL", 3: "EXTRA"}),
		"busy_timeout_ms":    pragmaInt(ctx, db, "busy_timeout"),
		"foreign_keys":       pragmaBool(ctx, db, "foreign_keys"),
		"temp_store":         pragmaEnum(ctx, db, "temp_store", map[int64]string{0: "DEFAULT", 1: "FILE", 2: "MEMORY"}),
		"journal_size_limit": pragmaInt(ctx, db, "journal_size_limit"),
		"wal_autocheckpoint": pragmaInt(ctx, db, "wal_autocheckpoint"),
	}

	rc := reqCount.Value()
	ec := errCount.Value()
	tl := totalLatencyMS.Value()

	epErrs := map[string]int64{}
	errCountByPath.Do(func(kv expvar.KeyValue) {
		if v, ok := kv.Value.(*expvar.Int); ok {
			epErrs[kv.Key] = v.Value()
		}
	})

	p95, p99 := latencyPercentiles()
	req1m, err1m := oneMinuteCounts()

	var avg, rate, rps, errRate1m float64
	if rc > 0 {
		avg = float64(tl) / float64(rc)
		rate = float64(ec) / float64(rc)
	}
	if req1m > 0 {
		rps = float64(req1m) / 60.0
		errRate1m = float64(err1m) / float64(req1m)
	}

	dbs := DBSnapshot{
		Open:    s.OpenConnections,
		MaxOpen: s.MaxOpenConnections,
		InUse:   s.InUse,
		Idle:    s.Idle,
		WaitCnt: s.WaitCount,
		WaitDur: s.WaitDuration,
		WaitRate: func() float64 {
			if rc > 0 {
				return float64(s.WaitCount) / float64(rc)
			}
			return 0
		}(),
		Pragma: pr,
	}

	http := HTTPSnapshot{
		AvgResMS:     avg,
		P95ResMS:     p95,
		P99ResMS:     p99,
		RPS:          rps,
		ErrorRate:    rate,
		ErrRate1m:    errRate1m,
		TotalReqs:    rc,
		TotalErrors:  ec,
		ErrorsByPath: epErrs,
	}

	return AppSnapshot{
		At:      time.Now(),
		DB:      dbs,
		Mem:     mem,
		Runtime: rt,
		HTTP:    http,
		Note:    "OK",
	}
}

func bytesToMB(b uint64) float64 {
	mb := float64(b) / 1024.0 / 1024.0
	return math.Round(mb*100) / 100 // 소수점 둘째자리까지만
}

func VarsUI(w http.ResponseWriter, r *http.Request) {
	service := os.Getenv("SERVICE_NAME")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, varsPage, service)
}

const varsPage = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Debug Vars</title>
  <style>
    body { font:14px sans-serif; margin:24px; }
    .card { border:1px solid #eee; padding:12px; margin-bottom:12px; }
    .row { display:flex; gap:8px; align-items:center; }
    button { padding:6px 10px; border:1px solid #ddd; background:#f8f8f8; cursor:pointer; }
  </style>
</head>
<body>
  <h1>Go Runtime / DB Snapshot</h1>
  <div class="row">
    <p style="margin:0"><a href="/debug/vars">/debug/vars</a></p>
    <button id="refreshBtn">새로고침</button>
    <span id="ts" style="color:#777; font-size:12px"></span>
  </div>
  <div id="app" style="margin-top:12px;"></div>

  <script>
    async function load() {
      const url = new URL(window.location.href);
      const auth = url.searchParams.get('auth');
      const r = await fetch("/debug/vars?auth=" + auth, { cache: 'no-store' });
      const d = await r.json();
      const a = d['%s'] || {};
      const m = a.mem || {};
      const b = a.db || {};
      const rt = a.runtime || {};
      const h = a.http || {};

      const avg = Number(h.avg_res_ms || 0).toFixed(2);
      const p95 = Number(h.p95_res_ms || 0).toFixed(2);
      const p99 = Number(h.p99_res_ms || 0).toFixed(2);
      const rps = Number(h.rps || 0).toFixed(2);
      const rate = Number((h.error_rate || 0) * 100).toFixed(2);
      const rate1m = Number((h.err_rate_1m || 0) * 100).toFixed(2);
      const waitRate = Number((b.wait_rate || 0) * 100).toFixed(2);

      const ep = h.errors_by_path || {};
      let epHtml = "";
      for (const k in ep) {
        epHtml += "<div>" + k + ": " + ep[k] + "</div>";
      }

      document.getElementById("app").innerHTML =
        "<div class='card'>" +
          "<b>메모리</b>" +
          "<div>Alloc: " + m.alloc_mb + " MB</div>" +
          "<div>Heap: " + m.heap_alloc_mb + "/" + m.heap_sys_mb + " MB</div>" +
          "<div>NumGC: " + m.num_gc + "</div>" +
          "<div>NextGC: " + m.next_gc_mb + " MB</div>" +
          "<div>PauseTotal: " + m.pause_total_ms + " ms</div>" +
        "</div>" +
        "<div class='card'>" +
          "<b>DB 연결</b>" +
          "<div>Open: " + b.open + "</div>" +
          "<div>InUse/Idle: " + b.in_use + "/" + b.idle + "</div>" +
          "<div>WaitCount: " + b.wait_count + "</div>" +
          "<div>WaitRate: " + waitRate + "%%</div>" +
        "</div>" +
        "<div class='card'>" +
          "<b>런타임</b>" +
          "<div>Goroutines: " + rt.goroutines + "</div>" +
        "</div>" +
        "<div class='card'>" +
          "<b>HTTP 요청</b>" +
          "<div>평균 응답시간: " + avg + " ms</div>" +
          "<div>P95/P99: " + p95 + "/" + p99 + " ms</div>" +
          "<div>RPS(1m): " + rps + "</div>" +
          "<div>전체 에러율: " + rate + "%% (" + (h.total_errors || 0) + "/" + (h.total_reqs || 0) + ")</div>" +
          "<div>최근1분 에러율: " + rate1m + "%%</div>" +
          "<div>엔드포인트별 에러:</div>" +
          epHtml +
        "</div>";

      var now = new Date();
      document.getElementById("ts").textContent = "업데이트: " + now.toLocaleString();
    }

    document.getElementById("refreshBtn").addEventListener("click", load);
    load(); // 최초 1회만
  </script>
</body>
</html>`

func pragmaString(ctx context.Context, db *sql.DB, name string) string {
	var v string
	_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
	return v
}
func pragmaInt(ctx context.Context, db *sql.DB, name string) int64 {
	var v int64
	_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
	return v
}
func pragmaBool(ctx context.Context, db *sql.DB, name string) bool {
	var v int64
	_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
	return v == 1
}
func pragmaEnum(ctx context.Context, db *sql.DB, name string, m map[int64]string) any {
	var v int64
	if err := db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v); err == nil {
		if s, ok := m[v]; ok {
			return s
		}
		return v
	}
	return nil
}
