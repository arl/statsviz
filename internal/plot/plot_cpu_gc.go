package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	name: "cpu-gc",
	tags: []tag{tagCPU, tagGC},
	metrics: []string{
		"/cpu/classes/gc/mark/assist:cpu-seconds",
		"/cpu/classes/gc/mark/dedicated:cpu-seconds",
		"/cpu/classes/gc/mark/idle:cpu-seconds",
		"/cpu/classes/gc/pause:cpu-seconds",
	},
	getvalues: func() getvalues {
		var (
			assist    = ratefloat64(idxcpuclassesgcmarkassist)
			dedicated = ratefloat64(idxcpuclassesgcmarkdedicated)
			idle      = ratefloat64(idxcpuclassesgcmarkidle)
			pause     = ratefloat64(idxcpuclassesgcpause)
		)

		return func(now time.Time, samples []metrics.Sample) any {
			return []float64{
				assist(now, samples),
				dedicated(now, samples),
				idle(now, samples),
				pause(now, samples),
			}
		}
	},
	layout: Scatter{
		Name:   "TODO(set later)",
		Title:  "CPU (Garbage Collector)",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "cpu-seconds per seconds",
				TickSuffix: "s",
			},
		},
		Subplots: []Subplot{
			{Name: "mark assist", Unitfmt: "%{y:.4s}s"},
			{Name: "mark dedicated", Unitfmt: "%{y:.4s}s"},
			{Name: "mark idle", Unitfmt: "%{y:.4s}s"},
			{Name: "pause", Unitfmt: "%{y:.4s}s"},
		},

		InfoText: `Cumulative metrics are converted to rates by Statsviz so as to be more easily comparable and readable.
All this metrics are overestimates, and not directly comparable to system CPU time measurements. Compare only with other /cpu/classes metrics.

<i>mark assist</i> is <b>/cpu/classes/gc/mark/assist</b>, estimated total CPU time goroutines spent performing GC tasks to assist the GC and prevent it from falling behind the application.
<i>mark dedicated</i> is <b>/cpu/classes/gc/mark/dedicated</b>, Estimated total CPU time spent performing GC tasks on processors (as defined by GOMAXPROCS) dedicated to those tasks.
<i>mark idle</i> is <b>/cpu/classes/gc/mark/idle</b>, estimated total CPU time spent performing GC tasks on spare CPU resources that the Go scheduler could not otherwise find a use for.
<i>pause</i> is <b>/cpu/classes/gc/pause</b>, estimated total CPU time spent with the application paused by the GC.

All metrics are rates in CPU-seconds per second.`,
	},
})
