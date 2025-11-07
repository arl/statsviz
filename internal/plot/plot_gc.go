package plot

import (
	"math"
	"runtime/metrics"
)

var _ = register(description{
	name: "garbage collection",
	tags: []tag{tagGC},
	metrics: []string{
		"/gc/gomemlimit:bytes",
		"/gc/heap/live:bytes",
		"/gc/heap/goal:bytes",
		"/memory/classes/total:bytes",
		"/memory/classes/heap/released:bytes",
	},
	layout: Scatter{
		Name:   "garbage collection",
		Title:  "GC Memory Summary",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "bytes",
				TickSuffix: "B",
			},
		},
		Subplots: []Subplot{
			{Name: "memory limit", Unitfmt: "%{y:.4s}B"},
			{Name: "in-use memory", Unitfmt: "%{y:.4s}B"},
			{Name: "heap live", Unitfmt: "%{y:.4s}B"},
			{Name: "heap goal", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `
<i>Memory limit</i> is <b>/gc/gomemlimit:bytes</b>, the Go runtime memory limit configured by the user (via GOMEMLIMIT or debug.SetMemoryLimt), otherwise 0. 
<i>In-use memory</i> is the total mapped memory minus released heap memory (<b>/memory/classes/total - /memory/classes/heap/released</b>).
<i>Heap live</i> is <b>/gc/heap/live:bytes</b>, heap memory occupied by live objects.  
<i>Heap goal</i> is <b>/gc/heap/goal:bytes</b>, the heap size target at the end of each GC cycle.`,
	},
	make: func(indices ...int) metricsGetter {
		return &garbageCollection{
			idxmemlimit:     indices[0],
			idxheaplive:     indices[1],
			idxheapgoal:     indices[2],
			idxmemtotal:     indices[3],
			idxheapreleased: indices[4],
		}
	},
})

type garbageCollection struct {
	idxmemlimit     int
	idxheaplive     int
	idxheapgoal     int
	idxmemtotal     int
	idxheapreleased int
}

func (p *garbageCollection) values(samples []metrics.Sample) any {
	memLimit := samples[p.idxmemlimit].Value.Uint64()
	heapLive := samples[p.idxheaplive].Value.Uint64()
	heapGoal := samples[p.idxheapgoal].Value.Uint64()
	memTotal := samples[p.idxmemtotal].Value.Uint64()
	heapReleased := samples[p.idxheapreleased].Value.Uint64()

	if memLimit == math.MaxInt64 {
		memLimit = 0
	}

	return []uint64{
		memLimit,
		memTotal - heapReleased,
		heapLive,
		heapGoal,
	}
}
