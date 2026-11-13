//go:build go1.26

package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/gc/finalizers/executed:finalizers",
		"/gc/finalizers/queued:finalizers",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			executed := samples[idx_gc_finalizers_executed_finalizers].Value.Uint64()
			queued := samples[idx_gc_finalizers_queued_finalizers].Value.Uint64()

			return []uint64{
				queued - executed,
			}
		}
	},
	layout: Scatter{
		Name:   "gc-finalizers",
		Tags:   []tag{tagGC},
		Title:  "GC Finalizers Queue",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			BarMode: "stack",
			Yaxis: ScatterYAxis{
				Title: "finalizers",
			},
		},
		Subplots: []Subplot{
			{Name: "queue size", Unitfmt: "%{y}"},
		},
		InfoText: `<i>queue size</i> is <b>/gc/finalizers/queued:finalizers - /gc/finalizers/executed:finalizers</b>.`,
	},
})
