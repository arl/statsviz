package plot

import (
	"runtime/metrics"
	"time"
)

type cpuOverall struct {
	idxUser     int
	idxScavenge int
	idxIdle     int
	idxGCtotal  int
	idxTotal    int

	lastTime     time.Time
	lastUser     float64
	lastScavenge float64
	lastIdle     float64
	lastGCtotal  float64
	lastTotal    float64
}

func makeCPUoverall(indices ...int) metricsGetter {
	return &cpuOverall{
		idxUser:     indices[0],
		idxScavenge: indices[1],
		idxIdle:     indices[2],
		idxGCtotal:  indices[3],
		idxTotal:    indices[4],
	}
}

func (p *cpuOverall) values(samples []metrics.Sample) any {
	curUser := samples[p.idxUser].Value.Float64()
	curScavenge := samples[p.idxScavenge].Value.Float64()
	curIdle := samples[p.idxIdle].Value.Float64()
	curGCtotal := samples[p.idxGCtotal].Value.Float64()
	curTotal := samples[p.idxTotal].Value.Float64()

	if p.lastTime.IsZero() {
		p.lastUser = curUser
		p.lastScavenge = curScavenge
		p.lastIdle = curIdle
		p.lastGCtotal = curGCtotal
		p.lastTotal = curTotal

		p.lastTime = time.Now()
		return []float64{0, 0, 0, 0, 0}
	}

	t := time.Since(p.lastTime).Seconds()

	user := (curUser - p.lastUser) / t
	scavenge := (curScavenge - p.lastScavenge) / t
	idle := (curIdle - p.lastIdle) / t
	gcTotal := (curGCtotal - p.lastGCtotal) / t
	total := (curTotal - p.lastTotal) / t

	p.lastUser = curUser
	p.lastScavenge = curScavenge
	p.lastIdle = curIdle
	p.lastGCtotal = curGCtotal
	p.lastTotal = curTotal

	return []float64{
		user,
		scavenge,
		idle,
		gcTotal,
		total,
	}
}

var cpuOverallLayout = Scatter{
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
}
