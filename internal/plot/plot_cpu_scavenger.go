package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/cpu/classes/scavenge/assist:cpu-seconds",
		"/cpu/classes/scavenge/background:cpu-seconds",
	},
	getvalues: func() getvalues {
		rateassist := rate[float64]()
		ratebackground := rate[float64]()

		return func(now time.Time, samples []metrics.Sample) any {
			assist := samples[idx_cpu_classes_scavenge_assist_cpu_seconds].Value.Float64()
			background := samples[idx_cpu_classes_scavenge_background_cpu_seconds].Value.Float64()

			return []float64{
				rateassist(now, assist),
				ratebackground(now, background),
			}
		}
	},
	layout: Scatter{
		Name:   "cpu-scavenger",
		Tags:   []tag{tagCPU, tagGC},
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
			{Name: "assist", Unitfmt: "%{y:.4s}s"},
			{Name: "background", Unitfmt: "%{y:.4s}s"},
		},
		InfoText: `<i>assist</i> is the rate of change, per second, of <b>/cpu/classes/scavenge/assist:cpu-seconds</b>. Estimated total CPU time spent returning unused memory to the underlying platform in response eagerly to memory pressure.
<i>background</i> is the rate of change, per second, of <b>/cpu/classes/scavenge/background:cpu-seconds</b>. Estimated total CPU time spent performing background tasks to return unused memory to the underlying platform.`,
	},
})
