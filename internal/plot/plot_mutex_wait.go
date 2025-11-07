package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	name: "mutex-wait",
	tags: []tag{tagMisc},
	metrics: []string{
		"/sync/mutex/wait/total:seconds",
	},
	layout: Scatter{
		Name:   "TODO(set later)",
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
	make: func(indices ...int) metricsGetter {
		return &mutexWait{
			idxMutexWait: indices[0],
		}
	},
})

type mutexWait struct {
	idxMutexWait int

	lastTime      time.Time
	lastMutexWait float64
}

func (p *mutexWait) values(samples []metrics.Sample) any {
	if p.lastTime.IsZero() {
		p.lastTime = time.Now()
		p.lastMutexWait = samples[p.idxMutexWait].Value.Float64()

		return []float64{0}
	}

	t := time.Since(p.lastTime).Seconds()

	mutexWait := (samples[p.idxMutexWait].Value.Float64() - p.lastMutexWait) / t

	p.lastMutexWait = samples[p.idxMutexWait].Value.Float64()
	p.lastTime = time.Now()

	return []float64{
		mutexWait,
	}
}
