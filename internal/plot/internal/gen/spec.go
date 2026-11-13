package main

// plots describes the runtime plots. Only what can't be deduced from
// runtime/metrics is specified here: names, which metrics are plotted together
// and how they're combined or transformed. Everything else (units, axes, hover
// format, info text) is deduced by the generator.
//
// A series Expr is a single metric or an expression of metrics. Transform is
// the name of a transform function in internal/plot/helpers.go: delta, rate,
// cpuRate, gcRate or noLimit.
//
// For heatmaps, the series name is what's counted in the buckets.
//
// The order of plots is their order in the UI. At runtime, a plot is skipped
// if one of its metrics is unknown to the Go version, or if a plot with the
// same name has already been registered.
var plots = []Plot{
	{
		Name: "cgo", Title: "CGO Calls", Tags: "misc", Type: "bar",
		Series: []Series{
			{"calls from go to c", "delta", "/cgo/go-to-c-calls:calls"},
		},
	},
	{
		Name: "cpu-gc", Title: "CPU (Garbage Collector)", Tags: "cpu gc", Type: "bar",
		Series: []Series{
			{"mark assist", "cpuRate", "/cpu/classes/gc/mark/assist:cpu-seconds"},
			{"mark dedicated", "cpuRate", "/cpu/classes/gc/mark/dedicated:cpu-seconds"},
			{"mark idle", "cpuRate", "/cpu/classes/gc/mark/idle:cpu-seconds"},
			{"pause", "cpuRate", "/cpu/classes/gc/pause:cpu-seconds"},
		},
	},
	{
		Name: "cpu-overall", Title: "CPU (Overall)", Tags: "cpu", Type: "bar",
		Series: []Series{
			{"user", "cpuRate", "/cpu/classes/user:cpu-seconds"},
			{"scavenge", "cpuRate", "/cpu/classes/scavenge/total:cpu-seconds"},
			{"idle", "cpuRate", "/cpu/classes/idle:cpu-seconds"},
			{"gc total", "cpuRate", "/cpu/classes/gc/total:cpu-seconds"},
		},
	},
	{
		Name: "cpu-scavenger", Title: "CPU (Scavenger)", Tags: "cpu gc", Type: "bar",
		Series: []Series{
			{"assist", "cpuRate", "/cpu/classes/scavenge/assist:cpu-seconds"},
			{"background", "cpuRate", "/cpu/classes/scavenge/background:cpu-seconds"},
		},
	},
	{
		Name: "garbage collection", Title: "GC Memory Summary", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"memory limit", "noLimit", "/gc/gomemlimit:bytes"},
			{"in-use memory", "", "/memory/classes/total:bytes - /memory/classes/heap/released:bytes"},
			{"heap live", "", "/gc/heap/live:bytes"},
			{"heap goal", "", "/gc/heap/goal:bytes"},
		},
	},
	{
		Name: "gc-cleanups", Title: "GC Cleanups Queue", Tags: "gc", Type: "bar",
		Series: []Series{
			{"queue size", "", "/gc/cleanups/queued:cleanups - /gc/cleanups/executed:cleanups"},
		},
	},
	{
		Name: "gc-cycles", Title: "Completed GC Cycles", Tags: "gc", Type: "bar",
		Series: []Series{
			{"automatic", "delta", "/gc/cycles/automatic:gc-cycles"},
			{"forced", "delta", "/gc/cycles/forced:gc-cycles"},
		},
	},
	{
		Name: "gc-finalizers", Title: "GC Finalizers Queue", Tags: "gc", Type: "bar",
		Series: []Series{
			{"queue size", "", "/gc/finalizers/queued:finalizers - /gc/finalizers/executed:finalizers"},
		},
	},
	{
		Name: "gc-scan", Title: "GC Scan", Tags: "gc", Type: "bar",
		Series: []Series{
			{"scannable globals", "", "/gc/scan/globals:bytes"},
			{"scannable heap", "", "/gc/scan/heap:bytes"},
			{"scanned stack", "", "/gc/scan/stack:bytes"},
		},
	},
	{
		Name: "gc-stack-size", Title: "Goroutines stack starting size", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"new goroutines stack size", "", "/gc/stack/starting-size:bytes"},
		},
	},
	{
		Name: "goroutines", Title: "Goroutines", Tags: "scheduler", Type: "scatter",
		Series: []Series{
			{"created", "delta", "/sched/goroutines-created:goroutines"},
			{"goroutines", "", "/sched/goroutines:goroutines"},
			{"not in Go", "", "/sched/goroutines/not-in-go:goroutines"},
			{"runnable", "", "/sched/goroutines/runnable:goroutines"},
			{"running", "", "/sched/goroutines/running:goroutines"},
			{"waiting", "", "/sched/goroutines/waiting:goroutines"},
		},
	},
	// Fallback for Go < 1.26.
	{
		Name: "goroutines", Title: "Goroutines", Tags: "scheduler", Type: "scatter",
		Series: []Series{
			{"goroutines", "", "/sched/goroutines:goroutines"},
		},
	},
	{
		Name: "alloc-free-rate", Title: "Heap Allocation & Free Rates", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"allocs/sec", "gcRate", "/gc/heap/allocs:objects"},
			{"frees/sec", "gcRate", "/gc/heap/frees:objects"},
		},
	},
	{
		Name: "heap (details)", Title: "Heap (details)", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"heap sys", "", "/memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/free:bytes + /memory/classes/heap/released:bytes"},
			{"heap objects", "", "/memory/classes/heap/objects:bytes"},
			{"heap stacks", "", "/memory/classes/heap/stacks:bytes"},
			{"heap goal", "", "/gc/heap/goal:bytes"},
		},
	},
	{
		Name: "live-bytes", Title: "Live Bytes in Heap", Tags: "gc", Type: "bar",
		Series: []Series{
			{"live bytes", "", "/gc/heap/allocs:bytes - /gc/heap/frees:bytes"},
		},
	},
	{
		Name: "live-objects", Title: "Live Objects in Heap", Tags: "gc", Type: "bar",
		Series: []Series{
			{"live objects", "", "/gc/heap/objects:objects"},
		},
	},
	{
		Name: "memory-classes", Title: "Memory classes", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"os stacks", "", "/memory/classes/os-stacks:bytes"},
			{"other", "", "/memory/classes/other:bytes"},
			{"profiling buckets", "", "/memory/classes/profiling/buckets:bytes"},
			{"total", "", "/memory/classes/total:bytes"},
		},
	},
	{
		Name: "mspan-mcache", Title: "MSpan/MCache", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"mspan in-use", "", "/memory/classes/metadata/mspan/inuse:bytes"},
			{"mspan free", "", "/memory/classes/metadata/mspan/free:bytes"},
			{"mcache in-use", "", "/memory/classes/metadata/mcache/inuse:bytes"},
			{"mcache free", "", "/memory/classes/metadata/mcache/free:bytes"},
		},
	},
	{
		Name: "mutex-wait", Title: "Mutex wait time", Tags: "misc", Type: "bar",
		Series: []Series{
			{"mutex wait", "rate", "/sync/mutex/wait/total:seconds"},
		},
	},
	{
		Name: "runnable-time", Title: "Time Goroutines Spend in 'Runnable' state", Tags: "scheduler", Type: "heatmap", Colorscale: "GreenShades",
		Series: []Series{
			{"goroutines", "", "/sched/latencies:seconds"},
		},
	},
	{Name: "size-classes", Custom: "sizeClasses"}, // see plot_size_classes.go
	{
		Name: "stopping-pauses-gc", Title: "Stop-the-world Stopping Latencies (GC)", Tags: "scheduler gc", Type: "heatmap", Colorscale: "GreenShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/stopping/gc:seconds"},
		},
	},
	{
		Name: "stopping-pauses-other", Title: "Stop-the-world Stopping Latencies (Other)", Tags: "scheduler", Type: "heatmap", Colorscale: "GreenShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/stopping/other:seconds"},
		},
	},
	{
		Name: "threads", Title: "Threads", Tags: "scheduler", Type: "scatter",
		Series: []Series{
			{"threads", "", "/sched/threads/total:threads"},
		},
	},
	{
		Name: "total-pauses-gc", Title: "Stop-the-world Pause Latencies (Total)", Tags: "scheduler gc", Type: "heatmap", Colorscale: "PinkShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/total/gc:seconds"},
		},
	},
	{
		Name: "total-pauses-other", Title: "Stop-the-world Pause Latencies (Other)", Tags: "scheduler", Type: "heatmap", Colorscale: "PinkShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/total/other:seconds"},
		},
	},
}
