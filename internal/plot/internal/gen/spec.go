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
var plots = []Plot{
	{
		File: "plot_cgo.go", Name: "cgo", Title: "CGO Calls", Tags: "misc", Type: "bar",
		Series: []Series{
			{"calls from go to c", "delta", "/cgo/go-to-c-calls:calls"},
		},
	},
	{
		File: "plot_cpu_gc.go", Name: "cpu-gc", Title: "CPU (Garbage Collector)", Tags: "cpu gc", Type: "bar",
		Series: []Series{
			{"mark assist", "cpuRate", "/cpu/classes/gc/mark/assist:cpu-seconds"},
			{"mark dedicated", "cpuRate", "/cpu/classes/gc/mark/dedicated:cpu-seconds"},
			{"mark idle", "cpuRate", "/cpu/classes/gc/mark/idle:cpu-seconds"},
			{"pause", "cpuRate", "/cpu/classes/gc/pause:cpu-seconds"},
		},
	},
	{
		File: "plot_cpu_overall.go", Name: "cpu-overall", Title: "CPU (Overall)", Tags: "cpu", Type: "bar",
		Series: []Series{
			{"user", "cpuRate", "/cpu/classes/user:cpu-seconds"},
			{"scavenge", "cpuRate", "/cpu/classes/scavenge/total:cpu-seconds"},
			{"idle", "cpuRate", "/cpu/classes/idle:cpu-seconds"},
			{"gc total", "cpuRate", "/cpu/classes/gc/total:cpu-seconds"},
		},
	},
	{
		File: "plot_cpu_scavenger.go", Name: "cpu-scavenger", Title: "CPU (Scavenger)", Tags: "cpu gc", Type: "bar",
		Series: []Series{
			{"assist", "cpuRate", "/cpu/classes/scavenge/assist:cpu-seconds"},
			{"background", "cpuRate", "/cpu/classes/scavenge/background:cpu-seconds"},
		},
	},
	{
		File: "plot_gc.go", Name: "garbage collection", Title: "GC Memory Summary", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"memory limit", "noLimit", "/gc/gomemlimit:bytes"},
			{"in-use memory", "", "/memory/classes/total:bytes - /memory/classes/heap/released:bytes"},
			{"heap live", "", "/gc/heap/live:bytes"},
			{"heap goal", "", "/gc/heap/goal:bytes"},
		},
	},
	{
		File: "plot_gc_cleanups.go", Build: "go1.26", Name: "gc-cleanups", Title: "GC Cleanups Queue", Tags: "gc", Type: "bar",
		Series: []Series{
			{"queue size", "", "/gc/cleanups/queued:cleanups - /gc/cleanups/executed:cleanups"},
		},
	},
	{
		File: "plot_gc_cycles.go", Name: "gc-cycles", Title: "Completed GC Cycles", Tags: "gc", Type: "bar",
		Series: []Series{
			{"automatic", "delta", "/gc/cycles/automatic:gc-cycles"},
			{"forced", "delta", "/gc/cycles/forced:gc-cycles"},
		},
	},
	{
		File: "plot_gc_finalizers.go", Build: "go1.26", Name: "gc-finalizers", Title: "GC Finalizers Queue", Tags: "gc", Type: "bar",
		Series: []Series{
			{"queue size", "", "/gc/finalizers/queued:finalizers - /gc/finalizers/executed:finalizers"},
		},
	},
	{
		File: "plot_gc_scan.go", Name: "gc-scan", Title: "GC Scan", Tags: "gc", Type: "bar",
		Series: []Series{
			{"scannable globals", "", "/gc/scan/globals:bytes"},
			{"scannable heap", "", "/gc/scan/heap:bytes"},
			{"scanned stack", "", "/gc/scan/stack:bytes"},
		},
	},
	{
		File: "plot_gc_stack_size.go", Name: "gc-stack-size", Title: "Goroutines stack starting size", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"new goroutines stack size", "", "/gc/stack/starting-size:bytes"},
		},
	},
	{
		File: "plot_goroutines.go", Build: "go1.26", Name: "goroutines", Title: "Goroutines", Tags: "scheduler", Type: "scatter",
		Series: []Series{
			{"created", "delta", "/sched/goroutines-created:goroutines"},
			{"goroutines", "", "/sched/goroutines:goroutines"},
			{"not in Go", "", "/sched/goroutines/not-in-go:goroutines"},
			{"runnable", "", "/sched/goroutines/runnable:goroutines"},
			{"running", "", "/sched/goroutines/running:goroutines"},
			{"waiting", "", "/sched/goroutines/waiting:goroutines"},
		},
	},
	{
		File: "plot_goroutines_pre_go1.26.go", Build: "!go1.26", Name: "goroutines", Title: "Goroutines", Tags: "scheduler", Type: "scatter",
		Series: []Series{
			{"goroutines", "", "/sched/goroutines:goroutines"},
		},
	},
	{
		File: "plot_heap_allocs.go", Name: "alloc-free-rate", Title: "Heap Allocation & Free Rates", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"allocs/sec", "gcRate", "/gc/heap/allocs:objects"},
			{"frees/sec", "gcRate", "/gc/heap/frees:objects"},
		},
	},
	{
		File: "plot_heap_details.go", Name: "heap (details)", Title: "Heap (details)", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"heap sys", "", "/memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/free:bytes + /memory/classes/heap/released:bytes"},
			{"heap objects", "", "/memory/classes/heap/objects:bytes"},
			{"heap stacks", "", "/memory/classes/heap/stacks:bytes"},
			{"heap goal", "", "/gc/heap/goal:bytes"},
		},
	},
	{
		File: "plot_live_bytes.go", Name: "live-bytes", Title: "Live Bytes in Heap", Tags: "gc", Type: "bar",
		Series: []Series{
			{"live bytes", "", "/gc/heap/allocs:bytes - /gc/heap/frees:bytes"},
		},
	},
	{
		File: "plot_live_objects.go", Name: "live-objects", Title: "Live Objects in Heap", Tags: "gc", Type: "bar",
		Series: []Series{
			{"live objects", "", "/gc/heap/objects:objects"},
		},
	},
	{
		File: "plot_memory_classes.go", Name: "memory-classes", Title: "Memory classes", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"os stacks", "", "/memory/classes/os-stacks:bytes"},
			{"other", "", "/memory/classes/other:bytes"},
			{"profiling buckets", "", "/memory/classes/profiling/buckets:bytes"},
			{"total", "", "/memory/classes/total:bytes"},
		},
	},
	{
		File: "plot_mspan_mcache.go", Name: "mspan-mcache", Title: "MSpan/MCache", Tags: "gc", Type: "scatter",
		Series: []Series{
			{"mspan in-use", "", "/memory/classes/metadata/mspan/inuse:bytes"},
			{"mspan free", "", "/memory/classes/metadata/mspan/free:bytes"},
			{"mcache in-use", "", "/memory/classes/metadata/mcache/inuse:bytes"},
			{"mcache free", "", "/memory/classes/metadata/mcache/free:bytes"},
		},
	},
	{
		File: "plot_mutex_wait.go", Name: "mutex-wait", Title: "Mutex wait time", Tags: "misc", Type: "bar",
		Series: []Series{
			{"mutex wait", "rate", "/sync/mutex/wait/total:seconds"},
		},
	},
	{
		File: "plot_runnable_time.go", Name: "runnable-time", Title: "Time Goroutines Spend in 'Runnable' state", Tags: "scheduler", Type: "heatmap", Colorscale: "GreenShades",
		Series: []Series{
			{"goroutines", "", "/sched/latencies:seconds"},
		},
	},
	{
		File: "plot_stopping_pauses_gc.go", Name: "stopping-pauses-gc", Title: "Stop-the-world Stopping Latencies (GC)", Tags: "scheduler gc", Type: "heatmap", Colorscale: "GreenShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/stopping/gc:seconds"},
		},
	},
	{
		File: "plot_stopping_pauses_other.go", Name: "stopping-pauses-other", Title: "Stop-the-world Stopping Latencies (Other)", Tags: "scheduler", Type: "heatmap", Colorscale: "GreenShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/stopping/other:seconds"},
		},
	},
	{
		File: "plot_threads.go", Build: "go1.26", Name: "threads", Title: "Threads", Tags: "scheduler", Type: "scatter",
		Series: []Series{
			{"threads", "", "/sched/threads/total:threads"},
		},
	},
	{
		File: "plot_total_pauses_gc.go", Name: "total-pauses-gc", Title: "Stop-the-world Pause Latencies (Total)", Tags: "scheduler gc", Type: "heatmap", Colorscale: "PinkShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/total/gc:seconds"},
		},
	},
	{
		File: "plot_total_pauses_other.go", Name: "total-pauses-other", Title: "Stop-the-world Pause Latencies (Other)", Tags: "scheduler", Type: "heatmap", Colorscale: "PinkShades",
		Series: []Series{
			{"pauses", "", "/sched/pauses/total/other:seconds"},
		},
	},
	// plot_size_classes.go is hand-written.
}
