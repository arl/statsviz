package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/cpu/classes/gc/mark/assist:cpu-seconds",
		"/cpu/classes/gc/mark/dedicated:cpu-seconds",
		"/cpu/classes/gc/mark/idle:cpu-seconds",
		"/cpu/classes/gc/pause:cpu-seconds",
	},
	getvalues: func() getvalues {
		rateassist := rate[float64]()
		ratededicated := rate[float64]()
		rateidle := rate[float64]()
		ratepause := rate[float64]()

		return func(now time.Time, samples []metrics.Sample) any {
			assist := samples[idx_cpu_classes_gc_mark_assist_cpu_seconds].Value.Float64()
			dedicated := samples[idx_cpu_classes_gc_mark_dedicated_cpu_seconds].Value.Float64()
			idle := samples[idx_cpu_classes_gc_mark_idle_cpu_seconds].Value.Float64()
			pause := samples[idx_cpu_classes_gc_pause_cpu_seconds].Value.Float64()

			return []float64{
				rateassist(now, assist),
				ratededicated(now, dedicated),
				rateidle(now, idle),
				ratepause(now, pause),
			}
		}
	},
	layout: Scatter{
		Name:   "cpu-gc",
		Tags:   []tag{tagCPU, tagGC},
		Title:  "CPU (Garbage Collector)",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "cpu-seconds / second",
				TickSuffix: "s",
			},
		},
		Subplots: []Subplot{
			{Name: "mark assist", Unitfmt: "%{y:.4s}s"},
			{Name: "mark dedicated", Unitfmt: "%{y:.4s}s"},
			{Name: "mark idle", Unitfmt: "%{y:.4s}s"},
			{Name: "pause", Unitfmt: "%{y:.4s}s"},
		},
		InfoText: `<i>mark assist</i> is the rate of change, per second, of <b>/cpu/classes/gc/mark/assist:cpu-seconds</b>. Estimated total CPU time goroutines spent performing GC tasks to assist the GC and prevent it from falling behind the application.
<i>mark dedicated</i> is the rate of change, per second, of <b>/cpu/classes/gc/mark/dedicated:cpu-seconds</b>. Estimated total CPU time spent performing GC tasks on processors (as defined by GOMAXPROCS) dedicated to those tasks.
<i>mark idle</i> is the rate of change, per second, of <b>/cpu/classes/gc/mark/idle:cpu-seconds</b>. Estimated total CPU time spent performing GC tasks on spare CPU resources that the Go scheduler could not otherwise find a use for.
<i>pause</i> is the rate of change, per second, of <b>/cpu/classes/gc/pause:cpu-seconds</b>. Estimated total CPU time spent with the application paused by the GC.`,
	},
})
