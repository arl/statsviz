package plot

import (
	"fmt"
	"math"
	"reflect"
	"runtime/metrics"
	"testing"
)

func Test_downsampleFactor(t *testing.T) {
	tests := []struct {
		nbuckets   int
		maxbuckets int
		want       int
	}{
		{nbuckets: 99, maxbuckets: 100, want: 1},
		{nbuckets: 100, maxbuckets: 100, want: 1},
		{nbuckets: 101, maxbuckets: 100, want: 2},
		{nbuckets: 10, maxbuckets: 5, want: 2},
		{nbuckets: 11, maxbuckets: 5, want: 3},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("n=%d,max=%d", tt.nbuckets, tt.maxbuckets), func(t *testing.T) {
			if got := downsampleFactor(tt.nbuckets, tt.maxbuckets); got != tt.want {
				t.Errorf("downsampleFactor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_downsample(t *testing.T) {
	tests := []struct {
		name        string
		hist        metrics.Float64Histogram
		factor      int
		wantBuckets []float64
		wantCounts  []uint64
	}{
		{
			name: "factor 1",
			hist: metrics.Float64Histogram{
				Buckets: []float64{0, 1, 2, 3, 4, 5, 6},
				Counts:  []uint64{2, 2, 1, 4, 5, 6},
			},
			factor:      1,
			wantBuckets: []float64{1, 2, 3, 4, 5, 6},
			wantCounts:  []uint64{2, 2, 1, 4, 5, 6},
		},
		{
			name: "factor 1 with infinites",
			hist: metrics.Float64Histogram{
				Buckets: []float64{math.Inf(-1), 1, 2, 3, 4, 5, math.Inf(1)},
				Counts:  []uint64{2, 2, 1, 4, 5, 6},
			},
			factor:      1,
			wantBuckets: []float64{1, 2, 3, 4, 5, 6},
			wantCounts:  []uint64{2, 2, 1, 4, 5, 6},
		},
		{
			name: "divisible by factor 3",
			hist: metrics.Float64Histogram{
				Buckets: []float64{0, 1, 2, 3, 4, 5, 6},
				Counts:  []uint64{2, 2, 1, 4, 5, 6},
			},
			factor:      3,
			wantBuckets: []float64{3, 6},
			wantCounts:  []uint64{5, 15},
		},
		{
			name: "divisible by factor 2",
			hist: metrics.Float64Histogram{
				Buckets: []float64{0, 1, 2, 3, 4, 5, 6},
				Counts:  []uint64{2, 2, 1, 4, 5, 6},
			},
			factor:      2,
			wantBuckets: []float64{2, 4, 6},
			wantCounts:  []uint64{4, 5, 11},
		},
		{
			name: "not divisible by factor 2",
			hist: metrics.Float64Histogram{
				Buckets: []float64{0, 1, 2, 3, 4, 5},
				Counts:  []uint64{2, 2, 1, 4, 5},
			},
			factor:      2,
			wantBuckets: []float64{2, 4, 5},
			wantCounts:  []uint64{4, 5, 5},
		},
		{
			name: "not divisible by factor 3, end +Inf",
			hist: metrics.Float64Histogram{
				Buckets: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, math.Inf(1)},
				Counts:  []uint64{2, 2, 1, 4, 5, 3, 2, 1, 7, 13},
			},
			factor:      3,
			wantBuckets: []float64{3, 6, 9, 12},
			wantCounts:  []uint64{5, 12, 10, 13},
		},
		{
			name: "not divisible by factor 3, start/end +Inf",
			hist: metrics.Float64Histogram{
				Buckets: []float64{math.Inf(-1), 1, 2, 3, 4, 5, 6, 7, 8, 9, math.Inf(1)},
				Counts:  []uint64{2, 2, 1, 4, 5, 3, 2, 1, 7, 13},
			},
			factor:      3,
			wantBuckets: []float64{3, 6, 9, 12},
			wantCounts:  []uint64{5, 12, 10, 13},
		},
		{
			name: "divisible by factor 3, end +Inf",
			hist: metrics.Float64Histogram{
				Buckets: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, math.Inf(1)},
				Counts:  []uint64{2, 2, 1, 4, 5, 3, 2, 1, 7},
			},
			factor:      3,
			wantBuckets: []float64{3, 6, 9},
			wantCounts:  []uint64{5, 12, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buckets := downsampleBuckets(&tt.hist, tt.factor)
			counts := downsampleCounts(&tt.hist, tt.factor)

			if !reflect.DeepEqual(buckets, tt.wantBuckets) {
				t.Errorf("downsampleBuckets() = %v, want %v", buckets, tt.wantBuckets)
			}
			if !reflect.DeepEqual(counts, tt.wantCounts) {
				t.Errorf("downsampleCounts() = %v, want %v", counts, tt.wantCounts)
			}
		})
	}
}
