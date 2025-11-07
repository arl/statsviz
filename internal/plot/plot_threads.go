package plot

import "runtime/metrics"

type threads struct {
	idxthreads int
}

func makeThreads(indices ...int) metricsGetter {
	return &threads{
		idxthreads: indices[0],
	}
}

func (p *threads) values(samples []metrics.Sample) any {
	threads := samples[p.idxthreads].Value.Uint64()
	return []uint64{
		threads,
	}
}

var threadsLayout = Scatter{
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
}
