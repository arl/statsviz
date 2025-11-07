package plot

import "runtime/metrics"

var _ = register(description{
	name: "threads",
	tags: []tag{tagScheduler},
	metrics: []string{
		"/sched/threads/total:threads",
	},
	layout: Scatter{
		Name:  "TODO(set later)",
		Title: "Threads",
		Type:  "scatter",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title: "bytes",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "threads",
				Unitfmt: "%{y}",
			},
		},
		InfoText: "Shows the current count of live threads that are owned by the Go runtime. Uses <b>/sched/threads/total:threads</b>",
	},
	make: func(idx ...int) metricsGetter {
		return &threads{
			idxthreads: idx[0],
		}
	},
})

type threads struct {
	idxthreads int
}

func (p *threads) values(samples []metrics.Sample) any {
	threads := samples[p.idxthreads].Value.Uint64()
	return []uint64{
		threads,
	}
}
