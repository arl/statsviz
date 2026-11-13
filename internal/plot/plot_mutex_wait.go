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
		ratemxwait := rate[float64]()

		return func(now time.Time, samples []metrics.Sample) any {
			mxwait := samples[idx_sync_mutex_wait_total_seconds].Value.Float64()

			return []float64{
				ratemxwait(now, mxwait),
			}
		}
	},
	layout: Scatter{
		Name:   "mutex-wait",
		Tags:   []tag{tagMisc},
		Title:  "Mutex wait time",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			BarMode: "stack",
			Yaxis: ScatterYAxis{
				Title:      "seconds / second",
				TickSuffix: "s",
			},
		},
		Subplots: []Subplot{
			{Name: "mutex wait", Unitfmt: "%{y:.4s}s"},
		},
		InfoText: `<i>mutex wait</i> is the rate of change, per second, of <b>/sync/mutex/wait/total:seconds</b>. Approximate cumulative time goroutines have spent blocked on a sync.Mutex, sync.RWMutex, or runtime-internal lock.`,
	},
})
