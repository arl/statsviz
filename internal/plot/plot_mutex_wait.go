package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/sync/mutex/wait/total:seconds",
	},
	getvalues: func() getvalues {
		ratemutexwait := rate[float64]()

		return func(now time.Time, samples []metrics.Sample) any {
			mutexwait := ratemutexwait(now, samples[idx_sync_mutex_wait_total_seconds].Value.Float64())

			return []float64{mutexwait}
		}
	},
	layout: Scatter{
		Name:   "mutex-wait",
		Tags:   []tag{tagMisc},
		Title:  "Mutex wait time",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "seconds / second",
				TickSuffix: "s",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "mutex wait",
				Unitfmt: "%{y:.4s}s",
				Type:    "bar",
			},
		},

		InfoText: `Cumulative metrics are converted to rates by Statsviz so as to be more easily comparable and readable.
<i>mutex wait</i> is <b>/sync/mutex/wait/total</b>, approximate cumulative time goroutines have spent blocked on a sync.Mutex or sync.RWMutex.

This metric is useful for identifying global changes in lock contention. Collect a mutex or block profile using the runtime/pprof package for more detailed contention data.`,
	},
})
