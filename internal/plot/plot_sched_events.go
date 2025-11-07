package plot

import (
	"math"
	"runtime/metrics"
)

var _ = register(description{
	name: "sched-events",
	tags: []tag{tagScheduler},
	metrics: []string{
		"/sched/latencies:seconds",
		"/sched/gomaxprocs:threads",
	},
	layout: Scatter{
		Name:   "TODO(set later)",
		Title:  "Goroutine Scheduling Events",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title: "events",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "events per unit of time",
				Unitfmt: "%{y}",
			},
			{
				Name:    "events per unit of time, per P",
				Unitfmt: "%{y}",
			},
		},
		InfoText: `<i>Events per second</i> is the sum of all buckets in <b>/sched/latencies:seconds</b>, that is, it tracks the total number of goroutine scheduling events. That number is multiplied by the constant 8.
<i>Events per second per P (processor)</i> is <i>Events per second</i> divided by current <b>GOMAXPROCS</b>, from <b>/sched/gomaxprocs:threads</b>.
<b>NOTE</b>: the multiplying factor comes from internal Go runtime source code and might change from version to version.`,
	},
	make: func(idx ...int) metricsGetter {
		return &schedEvents{
			idxschedlat:   idx[0],
			idxGomaxprocs: idx[1],
			lasttot:       math.MaxUint64,
		}
	},
})

type schedEvents struct {
	idxschedlat   int
	idxGomaxprocs int
	lasttot       uint64
}

// gTrackingPeriod is currently always 8. Guard it behind build tags when that
// changes. See https://github.com/golang/go/blob/go1.18.4/src/runtime/runtime2.go#L502-L504
const currentGtrackingPeriod = 8

// TODO show scheduling events per seconds
func (p *schedEvents) values(samples []metrics.Sample) any {
	schedlat := samples[p.idxschedlat].Value.Float64Histogram()
	gomaxprocs := samples[p.idxGomaxprocs].Value.Uint64()

	total := uint64(0)
	for _, v := range schedlat.Counts {
		total += v
	}
	total *= currentGtrackingPeriod

	curtot := total - p.lasttot
	if p.lasttot == math.MaxUint64 {
		// We don't want a big spike at statsviz launch in case the process has
		// been running for some time and curtot is high.
		curtot = 0
	}
	p.lasttot = total

	ftot := float64(curtot)

	return []float64{
		ftot,
		ftot / float64(gomaxprocs),
	}
}
