package debug

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestTakeSnapshotMetrics(t *testing.T) {
	db, err := sql.Open("sqlite", "file:memdb1?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("DB 열기 실패: %v", err)
	}
	defer db.Close()

	reqCount.Set(10)
	errCount.Set(2)
	totalLatencyMS.Set(500)
	t.Cleanup(func() {
		reqCount.Set(0)
		errCount.Set(0)
		totalLatencyMS.Set(0)
	})

	snap := takeSnapshot(db)
	if snap.HTTP.AvgResMS != 50 {
		t.Fatalf("평균 응답시간 예상값 50, 실제 %v", snap.HTTP.AvgResMS)
	}
	if snap.HTTP.ErrorRate != 0.2 {
		t.Fatalf("에러율 예상값 0.2, 실제 %v", snap.HTTP.ErrorRate)
	}
	if snap.Runtime.Goroutines == 0 {
		t.Fatal("고루틴 수는 0이면 안 됨")
	}
}
