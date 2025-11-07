package plot

import "runtime/metrics"

type heapDetails struct {
	idxobj      int
	idxunused   int
	idxfree     int
	idxreleased int
	idxstacks   int
	idxgoal     int
}

func makeHeapDetails(indices ...int) metricsGetter {
	return &heapDetails{
		idxobj:      indices[0],
		idxunused:   indices[1],
		idxfree:     indices[2],
		idxreleased: indices[3],
		idxstacks:   indices[4],
		idxgoal:     indices[5],
	}
}

func (p *heapDetails) values(samples []metrics.Sample) any {
	heapObjects := samples[p.idxobj].Value.Uint64()
	heapUnused := samples[p.idxunused].Value.Uint64()
	heapFree := samples[p.idxfree].Value.Uint64()
	heapReleased := samples[p.idxreleased].Value.Uint64()
	heapStacks := samples[p.idxstacks].Value.Uint64()
	nextGC := samples[p.idxgoal].Value.Uint64()

	heapIdle := heapReleased + heapFree
	heapInUse := heapObjects + heapUnused
	heapSys := heapInUse + heapIdle

	return []uint64{
		heapSys,
		heapObjects,
		heapStacks,
		nextGC,
	}
}

var heapDetailslLayout = Scatter{
	Name:   "TODO(set later)",
	Title:  "Heap (details)",
	Type:   "scatter",
	Events: "lastgc",
	Layout: ScatterLayout{
		Yaxis: ScatterYAxis{
			Title:      "bytes",
			TickSuffix: "B",
		},
	},
	Subplots: []Subplot{
		{
			Name:    "heap sys",
			Unitfmt: "%{y:.4s}B",
		},
		{
			Name:    "heap objects",
			Unitfmt: "%{y:.4s}B",
		},
		{
			Name:    "heap stacks",
			Unitfmt: "%{y:.4s}B",
		},
		{
			Name:    "heap goal",
			Unitfmt: "%{y:.4s}B",
		},
	},
	InfoText: `
<i>Heap</i> sys is <b>/memory/classes/heap/{objects + unused + released + free}</b>. It's an estimate of all the heap memory obtained from the OS.
<i>Heap objects</i> is <b>/memory/classes/heap/objects</b>, the memory occupied by live objects and dead objects that have not yet been marked free by the GC.
<i>Heap stacks</i> is <b>/memory/classes/heap/stacks</b>, the memory used for stack space.
<i>Heap goal</i> is <b>gc/heap/goal</b>, the heap size target for the end of the GC cycle.`,
}
