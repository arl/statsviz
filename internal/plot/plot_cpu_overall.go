package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	name: "cpu-overall",
	tags: []tag{tagCPU},
	metrics: []string{
		"/cpu/classes/user:cpu-seconds",
		"/cpu/classes/scavenge/total:cpu-seconds",
		"/cpu/classes/idle:cpu-seconds",
		"/cpu/classes/gc/total:cpu-seconds",
		"/cpu/classes/total:cpu-seconds",
	},
	getvalues: func() getvalues {
		var (
			user     = ratefloat64(idxcpuclassesuser)
			scavenge = ratefloat64(idxcpuclassesscavengtetotal)
			idle     = ratefloat64(idxcpuclassesidle)
			gctotal  = ratefloat64(idxcpuclassesgctotal)
			total    = ratefloat64(idxcpuclassestotal)
		)

		return func(now time.Time, samples []metrics.Sample) any {
			return []float64{
				user(now, samples),
				scavenge(now, samples),
				idle(now, samples),
				gctotal(now, samples),
				total(now, samples),
			}
		}
	},
	layout: Scatter{
		Name:   "TODO(set later)",
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
			{
				Name:    "user",
				Unitfmt: "%{y:.4s}s",
				Type:    "bar",
			},
			{
				Name:    "scavenge",
				Unitfmt: "%{y:.4s}s",
				Type:    "bar",
			},
			{
				Name:    "idle",
				Unitfmt: "%{y:.4s}s",
				Type:    "bar",
			},
			{
				Name:    "gc total",
				Unitfmt: "%{y:.4s}s",
				Type:    "bar",
			},
			{
				Name:    "total",
				Unitfmt: "%{y:.4s}s",
				Type:    "scatter",
			},
		},
		InfoText: `Shows the fraction of CPU spent in your code vs. runtime vs. wasted. Helps track overall utilization and potential headroom.
<i>user is</i> the rate of <b>/cpu/classes/user:cpu-seconds</b>, the CPU time spent running user Go code.
<i>scavenge is</i> the rate of <b>/cpu/classes/scavenge:cpu-seconds</b>, the CPU time spent performing tasks that return unused memory to the OS.
<i>idle is</i> the rate of <b>/cpu/classes/idle:cpu-seconds</b>, the CPU time spent performing GC tasks on spare CPU resources that the Go scheduler could not otherwise find a use for.
<i>gc total is</i> the rate of <b>/cpu/classes/gc/total:cpu-seconds</b>, the CPU time spent performing GC tasks (sum of all metrics in <b>/cpu/classes/gc</b>)
<i>total is</i> the rate of <b>/cpu/classes/total:cpu-seconds</b>, the available CPU time for user Go code or the Go runtime, as defined by GOMAXPROCS. In other words, GOMAXPROCS integrated over the wall-clock duration this process has been executing for.

All metrics are rates in CPU-seconds per second.`,
	},
})
