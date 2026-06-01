package calendar

import (
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	c, err := Convert(time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if c.Gregorian != "2026-06-01" {
		t.Errorf("gregorian = %s", c.Gregorian)
	}
	if c.HijriYear != 1447 || c.HijriMonth != 12 || c.HijriDay != 15 {
		t.Errorf("hijri = %d-%d-%d, want 1447-12-15", c.HijriYear, c.HijriMonth, c.HijriDay)
	}
	if c.HijriMonthName != "Dhu al-Hijjah" {
		t.Errorf("month name = %s", c.HijriMonthName)
	}
}
