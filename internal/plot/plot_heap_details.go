package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/memory/classes/heap/objects:bytes",
		"/memory/classes/heap/unused:bytes",
		"/memory/classes/heap/free:bytes",
		"/memory/classes/heap/released:bytes",
		"/memory/classes/heap/stacks:bytes",
		"/gc/heap/goal:bytes",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			heapObjects := samples[idx_memory_classes_heap_objects_bytes].Value.Uint64()
			heapUnused := samples[idx_memory_classes_heap_unused_bytes].Value.Uint64()
			heapFree := samples[idx_memory_classes_heap_free_bytes].Value.Uint64()
			heapReleased := samples[idx_memory_classes_heap_released_bytes].Value.Uint64()
			heapStacks := samples[idx_memory_classes_heap_stacks_bytes].Value.Uint64()
			heapGoal := samples[idx_gc_heap_goal_bytes].Value.Uint64()

			return []uint64{
				heapObjects + heapUnused + heapFree + heapReleased,
				heapObjects,
				heapStacks,
				heapGoal,
			}
		}
	},
	layout: Scatter{
		Name:   "heap (details)",
		Tags:   []tag{tagGC},
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
			{Name: "heap sys", Unitfmt: "%{y:.4s}B"},
			{Name: "heap objects", Unitfmt: "%{y:.4s}B"},
			{Name: "heap stacks", Unitfmt: "%{y:.4s}B"},
			{Name: "heap goal", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `<i>heap sys</i> is <b>/memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/free:bytes + /memory/classes/heap/released:bytes</b>.
<i>heap objects</i> is <b>/memory/classes/heap/objects:bytes</b>. Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.
<i>heap stacks</i> is <b>/memory/classes/heap/stacks:bytes</b>. Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
<i>heap goal</i> is <b>/gc/heap/goal:bytes</b>. Heap size target for the end of the GC cycle.`,
	},
})
