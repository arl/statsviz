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
	make: func(idx ...int) metricsGetter {
		return &cpuScavenger{
			idxScavengeAssist:     idx[0],
			idxScavengeBackground: idx[1],
		}
	},
})

type cpuScavenger struct {
	idxScavengeAssist     int
	idxScavengeBackground int

	lastTime time.Time

	lastScavengeAssist     float64
	lastScavengeBackground float64
}

func (p *cpuScavenger) values(samples []metrics.Sample) any {
	curScavengeAssist := samples[p.idxScavengeAssist].Value.Float64()
	curScavengeBackground := samples[p.idxScavengeBackground].Value.Float64()

	if p.lastTime.IsZero() {
		p.lastScavengeAssist = curScavengeAssist
		p.lastScavengeBackground = curScavengeBackground
		p.lastTime = time.Now()

		return []float64{0, 0, 0, 0, 0}
	}

	t := time.Since(p.lastTime).Seconds()

	scavengeAssist := (curScavengeAssist - p.lastScavengeAssist) / t
	scavengeBackground := (curScavengeBackground - p.lastScavengeBackground) / t

	p.lastScavengeAssist = curScavengeAssist
	p.lastScavengeBackground = curScavengeBackground

	return []float64{
		scavengeAssist,
		scavengeBackground,
	}
}
