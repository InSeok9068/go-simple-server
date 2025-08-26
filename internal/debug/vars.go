package debug

import (
	"context"
	"database/sql"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

type DBSnapshot struct {
	Open    int           `json:"open"`
	MaxOpen int           `json:"max_open"`
	InUse   int           `json:"in_use"`
	Idle    int           `json:"idle"`
	WaitCnt int64         `json:"wait_count"`
	WaitDur time.Duration `json:"wait_duration"`
	Pragma  any           `json:"pragma"`
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

type AppSnapshot struct {
	At   time.Time   `json:"at"`
	DB   DBSnapshot  `json:"db"`
	Mem  MemSnapshot `json:"mem"`
	Note string      `json:"note"`
}

var snap atomic.Value

func Init(name string, db *sql.DB) {
	ctx := context.Background()
	snap.Store(takeSnapshot(ctx, db))
	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for range t.C {
			snap.Store(takeSnapshot(ctx, db))
		}
	}()
	expvar.Publish(name, expvar.Func(snap.Load))
}

func takeSnapshot(ctx context.Context, db *sql.DB) AppSnapshot {
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

	s := db.Stats()
	dbs := DBSnapshot{
		Open:    s.OpenConnections,
		MaxOpen: s.MaxOpenConnections,
		InUse:   s.InUse,
		Idle:    s.Idle,
		WaitCnt: s.WaitCount,
		WaitDur: s.WaitDuration,
		Pragma:  readSQLitePragmas(ctx, db),
	}

	return AppSnapshot{
		At:   time.Now(),
		DB:   dbs,
		Mem:  mem,
		Note: "OK",
	}
}

func bytesToMB(b uint64) float64 { return float64(b) / 1024.0 / 1024.0 }

func readSQLitePragmas(ctx context.Context, db *sql.DB) map[string]any {
	pragma := map[string]any{}
	readStr := func(name string) string {
		var v string
		_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
		return v
	}
	readInt := func(name string) int64 {
		var v int64
		_ = db.QueryRowContext(ctx, "PRAGMA "+name).Scan(&v)
		return v
	}

	pragma["busy_timeout_ms"] = readInt("busy_timeout")
	pragma["foreign_keys"] = readInt("foreign_keys") == 1
	pragma["journal_mode"] = readStr("journal_mode")
	pragma["journal_size_limit"] = readInt("journal_size_limit")
	pragma["synchronous"] = readStr("synchronous")
	pragma["temp_store"] = readStr("temp_store")
	pragma["wal_autocheckpoint"] = readInt("wal_autocheckpoint")
	return pragma
}

func VarsUI(w http.ResponseWriter, r *http.Request) {
	service := os.Getenv("SERVICE_NAME")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, varsPage, service)
}

const varsPage = "<!doctype html><html><head><meta charset=\"utf-8\"><title>Debug Vars</title>" +
	"<style>body{font:14px sans-serif;margin:24px}.card{border:1px solid #eee;padding:12px;margin-bottom:12px}</style>" +
	"</head><body><h1>Go Runtime / DB Snapshot</h1><p><a href=\"/debug/vars\">/debug/vars</a></p><div id=\"app\"></div>" +
	"<script>async function load(){const r=await fetch('/debug/vars');const d=await r.json();const a=d['%s']||{};const m=a.mem||{};const b=a.db||{};const f=v=>v==null?'-':v;document.getElementById('app').innerHTML=`<div class='card'><b>메모리</b><div>Alloc: ${m.alloc_mb} MB</div><div>Heap: ${m.heap_alloc_mb}/${m.heap_sys_mb} MB</div><div>NumGC: ${m.num_gc}</div></div><div class='card'><b>DB 연결</b><div>Open: ${b.open}</div><div>InUse/Idle: ${b.in_use}/${b.idle}</div><div>WaitCount: ${b.wait_count}</div></div>`;}load();setInterval(load,5000);</script></body></html>"
