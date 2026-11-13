package plot

import (
	"testing"
	"time"
)

func TestCPURate(t *testing.T) {
	rate := cpuRate()

	// total grows by GOMAXPROCS (4) CPU-seconds per second.
	steps := []struct {
		total, cur, want float64
	}{
		{0, 0, 0},    // no GC cycle yet
		{40, 10, 1},  // first update: average since program start
		{40, 10, 1},  // no update: rate is held
		{80, 30, 2},  // 20 CPU-seconds in 10s
		{80, 30, 2},  // held
		{120, 30, 0}, // nothing consumed in 10s
	}
	for i, s := range steps {
		if got := rate(s.total, 4, s.cur); got != s.want {
			t.Errorf("step %d: rate = %v, want %v", i, got, s.want)
		}
	}
}

func TestGCRate(t *testing.T) {
	rate := gcRate[uint64]()
	t0 := time.Now()
	at := func(sec int) time.Time { return t0.Add(time.Duration(sec) * time.Second) }

	steps := []struct {
		sec    int
		cycles uint64
		cur    uint64
		want   float64
	}{
		{0, 1, 100, 0},    // first sample
		{1, 1, 150, 0},    // same GC cycle: nothing to report yet
		{2, 2, 300, 100},  // 200 in 2s
		{3, 2, 300, 100},  // held
		{4, 2, 350, 100},  // held
		{6, 3, 700, 100},  // 400 in 4s
		{7, 5, 1000, 300}, // 300 in 1s
	}
	for i, s := range steps {
		if got := rate(at(s.sec), s.cycles, s.cur); got != s.want {
			t.Errorf("step %d: rate = %v, want %v", i, got, s.want)
		}
	}
}
