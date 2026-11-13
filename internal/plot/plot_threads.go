//go:build go1.26

package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/sched/threads/total:threads",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			threads := samples[idx_sched_threads_total_threads].Value.Uint64()

			return []uint64{
				threads,
			}
		}
	},
	layout: Scatter{
		Name:   "threads",
		Tags:   []tag{tagScheduler},
		Title:  "Threads",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:   "threads",
				Integer: true,
			},
		},
		Subplots: []Subplot{
			{Name: "threads", Unitfmt: "%{y}"},
		},
		InfoText: `<i>threads</i> is <b>/sched/threads/total:threads</b>. The current count of live threads that are owned by the Go runtime.`,
	},
})
