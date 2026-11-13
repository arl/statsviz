package plot

import "testing"

func TestHeldRate(t *testing.T) {
	rate := heldRate()

	steps := []struct {
		version, clock, cur, want float64
	}{
		{1, 0, 100, 0},   // first sample
		{1, 1, 150, 0},   // same version: nothing to report yet
		{2, 2, 300, 100}, // 200 in 2s
		{2, 3, 300, 100}, // held
		{2, 4, 350, 100}, // held
		{3, 6, 700, 100}, // 400 in 4s
		{5, 7, 1000, 300},
	}
	for i, s := range steps {
		if got := rate(s.version, s.clock, s.cur); got != s.want {
			t.Errorf("step %d: rate = %v, want %v", i, got, s.want)
		}
	}
}
