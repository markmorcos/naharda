package domain

import "testing"

func TestCities(t *testing.T) {
	if len(Cities) != 13 {
		t.Errorf("want 13 canonical cities, got %d", len(Cities))
	}
	if c, ok := CityBySlug("cairo"); !ok || c.Name != "Cairo" {
		t.Errorf("cairo lookup failed: %+v ok=%v", c, ok)
	}
	if _, ok := CityBySlug("atlantis"); ok {
		t.Error("atlantis should not be a known city")
	}
}
