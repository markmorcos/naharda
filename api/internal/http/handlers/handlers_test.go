package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/markmorcos/naharda/api/internal/store"
)

// noDB returns handlers backed by a store with no pool — exercises the paths
// that don't require a database (cities, calendar, fuel defaults, health).
func noDB(t *testing.T) *Handlers {
	t.Helper()
	st, err := store.New(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	return New(st, false, nil)
}

func TestCitiesHandler(t *testing.T) {
	w := httptest.NewRecorder()
	noDB(t).Cities(w, httptest.NewRequest(http.MethodGet, "/v1/cities", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "cairo") {
		t.Error("cities response missing cairo")
	}
}

func TestCalendarHandler(t *testing.T) {
	w := httptest.NewRecorder()
	noDB(t).Calendar(w, httptest.NewRequest(http.MethodGet, "/v1/calendar", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "hijri") {
		t.Error("calendar response missing hijri")
	}
}

func TestFuelHandlerUsesDefaults(t *testing.T) {
	w := httptest.NewRecorder()
	noDB(t).Fuel(w, httptest.NewRequest(http.MethodGet, "/v1/fuel", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "gasoline-92") {
		t.Error("fuel response missing gasoline-92")
	}
}

func TestHealthAndReadiness(t *testing.T) {
	h := noDB(t)
	w := httptest.NewRecorder()
	h.Healthz(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != http.StatusOK {
		t.Errorf("healthz = %d", w.Code)
	}
	w = httptest.NewRecorder()
	h.Readyz(w, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("readyz without DB = %d, want 503", w.Code)
	}
}

func TestAggregateParallel(t *testing.T) {
	agg := aggregateParallel([]store.FXRate{
		{Value: 61.5, Source: "a"},
		{Value: 62.8, Source: "b"},
		{Value: 62.1, Source: "c"},
	})
	if agg["n"].(int) != 3 {
		t.Errorf("n = %v", agg["n"])
	}
	if agg["min"].(float64) != 61.5 || agg["max"].(float64) != 62.8 {
		t.Errorf("min/max = %v/%v", agg["min"], agg["max"])
	}
	if agg["avg"].(float64) != 62.13 {
		t.Errorf("avg = %v, want 62.13", agg["avg"])
	}
}

func TestAggregateParallelEmpty(t *testing.T) {
	agg := aggregateParallel(nil)
	if agg["n"].(int) != 0 {
		t.Errorf("empty n = %v", agg["n"])
	}
	if agg["min"] != nil {
		t.Errorf("empty min = %v, want nil", agg["min"])
	}
}
