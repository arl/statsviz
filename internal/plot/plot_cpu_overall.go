package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/cpu/classes/user:cpu-seconds",
		"/cpu/classes/scavenge/total:cpu-seconds",
		"/cpu/classes/idle:cpu-seconds",
		"/cpu/classes/gc/total:cpu-seconds",
		"/cpu/classes/total:cpu-seconds",
	},
	getvalues: func() getvalues {
		rateuser := rate[float64]()
		ratescavenge := rate[float64]()
		rateidle := rate[float64]()
		rategctotal := rate[float64]()
		ratetotal := rate[float64]()

		return func(now time.Time, samples []metrics.Sample) any {
			user := samples[idx_cpu_classes_user_cpu_seconds].Value.Float64()
			scavenge := samples[idx_cpu_classes_scavenge_total_cpu_seconds].Value.Float64()
			idle := samples[idx_cpu_classes_idle_cpu_seconds].Value.Float64()
			gctotal := samples[idx_cpu_classes_gc_total_cpu_seconds].Value.Float64()
			total := samples[idx_cpu_classes_total_cpu_seconds].Value.Float64()

			return []float64{
				rateuser(now, user),
				ratescavenge(now, scavenge),
				rateidle(now, idle),
				rategctotal(now, gctotal),
				ratetotal(now, total),
			}
		}
	},
	layout: Scatter{
		Name:   "cpu-overall",
		Tags:   []tag{tagCPU},
		Title:  "CPU (Overall)",
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
			{Name: "user", Unitfmt: "%{y:.4s}s"},
			{Name: "scavenge", Unitfmt: "%{y:.4s}s"},
			{Name: "idle", Unitfmt: "%{y:.4s}s"},
			{Name: "gc total", Unitfmt: "%{y:.4s}s"},
			{Name: "total", Unitfmt: "%{y:.4s}s", Type: "scatter"},
		},
		InfoText: `<i>user</i> is the rate of change, per second, of <b>/cpu/classes/user:cpu-seconds</b>. Estimated total CPU time spent running user Go code.
<i>scavenge</i> is the rate of change, per second, of <b>/cpu/classes/scavenge/total:cpu-seconds</b>. Estimated total CPU time spent performing tasks that return unused memory to the underlying platform.
<i>idle</i> is the rate of change, per second, of <b>/cpu/classes/idle:cpu-seconds</b>. Estimated total available CPU time not spent executing any Go or Go runtime code.
<i>gc total</i> is the rate of change, per second, of <b>/cpu/classes/gc/total:cpu-seconds</b>. Estimated total CPU time spent performing GC tasks.
<i>total</i> is the rate of change, per second, of <b>/cpu/classes/total:cpu-seconds</b>. Estimated total available CPU time for user Go code or the Go runtime, as defined by GOMAXPROCS.`,
	},
})
