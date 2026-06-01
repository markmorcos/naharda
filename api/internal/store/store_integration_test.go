package store

import (
	"context"
	"os"
	"testing"

	"github.com/markmorcos/naharda/api/migrations"
)

// testStore connects to TEST_DATABASE_URL, runs migrations, and returns a Store.
// Skips when the env var is unset (e.g. local unit-only runs).
func testStore(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping integration test")
	}
	if err := Migrate(dsn, migrations.FS); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	st, err := New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(st.Close)
	return st
}

func TestFXRoundTrip(t *testing.T) {
	st := testStore(t)
	ctx := context.Background()
	if err := st.InsertFXRate(ctx, "official", "USD", 50.5, "test-src", false); err != nil {
		t.Fatal(err)
	}
	rates, err := st.LatestFXRates(ctx, "official")
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, r := range rates {
		if r.Quote == "USD" {
			found = true
		}
	}
	if !found {
		t.Error("USD not returned by LatestFXRates")
	}
	v, ok, err := st.LatestFXRate(ctx, "official", "USD")
	if err != nil || !ok || v <= 0 {
		t.Errorf("LatestFXRate USD = %v ok=%v err=%v", v, ok, err)
	}
}

func TestOutlierHeldExcludedFromLatest(t *testing.T) {
	st := testStore(t)
	ctx := context.Background()
	// A pending_review row must not be served as last-good.
	_ = st.InsertFXRate(ctx, "official", "EUR", 60.0, "test-src", false)
	_ = st.InsertFXRate(ctx, "official", "EUR", 999.0, "test-src", true) // held outlier
	v, ok, err := st.LatestFXRate(ctx, "official", "EUR")
	if err != nil || !ok {
		t.Fatalf("err=%v ok=%v", err, ok)
	}
	if v == 999.0 {
		t.Error("held outlier was served as latest")
	}
}

func TestGoldRoundTrip(t *testing.T) {
	st := testStore(t)
	ctx := context.Background()
	for _, k := range []int{18, 21, 24} {
		if err := st.InsertGoldPrice(ctx, "world_derived", k, float64(k)*100, "test", false); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := st.LatestGoldPrices(ctx, "world_derived")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) < 3 {
		t.Errorf("want >=3 karats, got %d", len(rows))
	}
}

func TestSignupAndStats(t *testing.T) {
	st := testStore(t)
	ctx := context.Background()
	email := "itest@example.com"
	if err := st.InsertSignup(ctx, email, true, "en"); err != nil {
		t.Fatal(err)
	}
	// Idempotent on email.
	if err := st.InsertSignup(ctx, email, true, "en"); err != nil {
		t.Fatal(err)
	}
	stats, err := st.GetStats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if stats.SignupsCount < 1 {
		t.Errorf("signups count = %d", stats.SignupsCount)
	}
}

func TestActiveOverride(t *testing.T) {
	st := testStore(t)
	ctx := context.Background()
	_, err := st.Pool.Exec(ctx,
		`INSERT INTO manual_override (family, key, value, effective_from) VALUES ('fuel','gasoline-92',13.99,now() - interval '1 hour')`)
	if err != nil {
		t.Fatal(err)
	}
	v, ok, err := st.ActiveOverride(ctx, "fuel", "gasoline-92")
	if err != nil || !ok {
		t.Fatalf("err=%v ok=%v", err, ok)
	}
	if v != 13.99 {
		t.Errorf("override value = %v", v)
	}
	// Expired override is not active.
	_, _ = st.Pool.Exec(ctx,
		`INSERT INTO manual_override (family, key, value, effective_from, effective_to) VALUES ('fuel','diesel',1,now() - interval '2 day', now() - interval '1 day')`)
	if _, ok, _ := st.ActiveOverride(ctx, "fuel", "diesel"); ok {
		t.Error("expired override should not be active")
	}
}
