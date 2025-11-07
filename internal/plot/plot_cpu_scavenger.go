package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	name: "cpu-scavenger",
	tags: []tag{tagCPU, tagGC},
	metrics: []string{
		"/cpu/classes/scavenge/assist:cpu-seconds",
		"/cpu/classes/scavenge/background:cpu-seconds",
	},
	getvalues: func() getvalues {
		var (
			assist     = ratefloat64(idxcpuclassesscavengeassist)
			background = ratefloat64(idxcpuclassesscavengebackground)
		)

		return func(now time.Time, samples []metrics.Sample) any {
			return []float64{
				assist(now, samples),
				background(now, samples),
			}
		}
	},
	layout: Scatter{
		Name:   "TODO(set later)",
		Title:  "CPU (Scavenger)",
		Type:   "bar",
		Events: "lastgc",
		Layout: ScatterLayout{
			BarMode: "stack",
			Yaxis: ScatterYAxis{
				Title:      "cpu-seconds / second",
				TickSuffix: "s",
			},
		},
		Subplots: []Subplot{
			{
				Name:    "assist",
				Unitfmt: "%{y:.4s}s",
				Type:    "bar",
			},
			{
				Name:    "background",
				Unitfmt: "%{y:.4s}s",
				Type:    "bar",
			},
		},
		InfoText: `Breakdown of how the GC scavenger returns memory to the OS (eagerly vs background).
<i>assist is</i> the rate of <b>/cpu/classes/scavenge/assist</b>, the CPU time spent returning unused memory eagerly in response to memory pressure.
<i>background is</i> the rate of <b>/cpu/classes/scavenge/background</b>, the CPU time spent performing background tasks to return unused memory to the OS.

Both metrics are rates in CPU-seconds per second.`,
	},
})
