package quality

import "testing"

func TestIsOutlier(t *testing.T) {
	cases := []struct {
		name            string
		value, avg, thr float64
		want            bool
	}{
		{"equal", 100, 100, 5, false},
		{"just under threshold", 104.9, 100, 5, false},
		{"over threshold", 106, 100, 5, true},
		{"negative deviation over", 90, 100, 5, true},
		{"no baseline", 100, 0, 5, false},
	}
	for _, c := range cases {
		if got := IsOutlier(c.value, c.avg, c.thr); got != c.want {
			t.Errorf("%s: IsOutlier(%v,%v,%v)=%v want %v", c.name, c.value, c.avg, c.thr, got, c.want)
		}
	}
}
