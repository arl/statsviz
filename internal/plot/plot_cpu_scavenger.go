package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/cpu/classes/scavenge/assist:cpu-seconds",
		"/cpu/classes/scavenge/background:cpu-seconds",
		"/cpu/classes/total:cpu-seconds",
		"/sched/gomaxprocs:threads",
	},
	getvalues: func() getvalues {
		rateassist := cpuRate()
		ratebackground := cpuRate()

		return func(_ time.Time, samples []metrics.Sample) any {
			assist := samples[idx_cpu_classes_scavenge_assist_cpu_seconds].Value.Float64()
			background := samples[idx_cpu_classes_scavenge_background_cpu_seconds].Value.Float64()
			total := samples[idx_cpu_classes_total_cpu_seconds].Value.Float64()
			maxprocs := samples[idx_sched_gomaxprocs_threads].Value.Uint64()

			return []float64{
				rateassist(total, maxprocs, assist),
				ratebackground(total, maxprocs, background),
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
		InfoText: `<i>assist</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/cpu/classes/scavenge/assist:cpu-seconds</b>. Estimated total CPU time spent returning unused memory to the underlying platform in response eagerly to memory pressure.
<i>background</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/cpu/classes/scavenge/background:cpu-seconds</b>. Estimated total CPU time spent performing background tasks to return unused memory to the underlying platform.`,
	},
})
