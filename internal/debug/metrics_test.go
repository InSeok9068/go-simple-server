package debug

import (
	"expvar"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMetricsMiddleware_ErrorCountByPath(t *testing.T) {
	errCountByPath.Init()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/test")

	h := func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "boom")
	}

	if err := MetricsMiddleware(h)(c); err == nil {
		t.Fatalf("expected error")
	}

	v := errCountByPath.Get("/test")
	if v == nil {
		t.Fatalf("count not recorded")
	}
	if count := v.(*expvar.Int).Value(); count != 1 {
		t.Fatalf("unexpected count: %d", count)
	}
}
