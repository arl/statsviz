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
		"/sched/gomaxprocs:threads",
	},
	getvalues: func() getvalues {
		rateuser := cpuRate()
		ratescavenge := cpuRate()
		rateidle := cpuRate()
		rategctotal := cpuRate()

		return func(_ time.Time, samples []metrics.Sample) any {
			user := samples[idx_cpu_classes_user_cpu_seconds].Value.Float64()
			scavenge := samples[idx_cpu_classes_scavenge_total_cpu_seconds].Value.Float64()
			idle := samples[idx_cpu_classes_idle_cpu_seconds].Value.Float64()
			gctotal := samples[idx_cpu_classes_gc_total_cpu_seconds].Value.Float64()
			total := samples[idx_cpu_classes_total_cpu_seconds].Value.Float64()
			maxprocs := samples[idx_sched_gomaxprocs_threads].Value.Uint64()

			return []float64{
				rateuser(total, maxprocs, user),
				ratescavenge(total, maxprocs, scavenge),
				rateidle(total, maxprocs, idle),
				rategctotal(total, maxprocs, gctotal),
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
		},
		InfoText: `<i>user</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/cpu/classes/user:cpu-seconds</b>. Estimated total CPU time spent running user Go code.
<i>scavenge</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/cpu/classes/scavenge/total:cpu-seconds</b>. Estimated total CPU time spent performing tasks that return unused memory to the underlying platform.
<i>idle</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/cpu/classes/idle:cpu-seconds</b>. Estimated total available CPU time not spent executing any Go or Go runtime code.
<i>gc total</i> is the rate of change, per second, averaged over the last GC cycle, of <b>/cpu/classes/gc/total:cpu-seconds</b>. Estimated total CPU time spent performing GC tasks.`,
	},
})
