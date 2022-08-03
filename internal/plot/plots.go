package plot

import "runtime/metrics"

/*
 * heap (global)
 */

type heapGlobal struct {
	enabled bool

	idxobj      int
	idxunused   int
	idxfree     int
	idxreleased int
}

func makeHeapGlobalPlot(am *allMetrics) *heapGlobal {
	idxobj, ok1 := am.idxs["/memory/classes/heap/objects:bytes"]
	idxunused, ok2 := am.idxs["/memory/classes/heap/unused:bytes"]
	idxfree, ok3 := am.idxs["/memory/classes/heap/free:bytes"]
	idxreleased, ok4 := am.idxs["/memory/classes/heap/released:bytes"]

	return &heapGlobal{
		enabled:     ok1 && ok2 && ok3 && ok4,
		idxobj:      idxobj,
		idxunused:   idxunused,
		idxfree:     idxfree,
		idxreleased: idxreleased,
	}
}

func (p *heapGlobal) name() string    { return "heap-global" }
func (p *heapGlobal) isEnabled() bool { return p.enabled }

func (p *heapGlobal) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:   p.name(),
		Title:  "Heap (global)",
		Type:   "scatter",
		Events: "lastgc",
		Subplots: []Subplot{
			{
				Name:       "heap in-use",
				Unitfmt:    "%{y:.4s}B",
				HoverOn:    "points+fills",
				StackGroup: "one",
			},
			{
				Name:       "heap free",
				Unitfmt:    "%{y:.4s}B",
				HoverOn:    "points+fills",
				StackGroup: "one",
			},
			{
				Name:       "heap released",
				Unitfmt:    "%{y:.4s}B",
				HoverOn:    "points+fills",
				StackGroup: "one",
			},
		},
	}
	s.Layout.Yaxis.TickSuffix = "B"
	s.Layout.Yaxis.Title = "bytes"
	return s
}

func (p *heapGlobal) values(samples []metrics.Sample) interface{} {
	heapObjects := samples[p.idxobj].Value.Uint64()
	heapUnused := samples[p.idxunused].Value.Uint64()

	heapInUse := heapObjects + heapUnused
	heapFree := samples[p.idxfree].Value.Uint64()
	heapReleased := samples[p.idxreleased].Value.Uint64()
	return []uint64{
		heapInUse,
		heapFree,
		heapReleased,
	}
}

/*
 * heap (details)
 */

type heapDetails struct {
	enabled bool

	idxobj      int
	idxunused   int
	idxfree     int
	idxreleased int
	idxstacks   int
	idxgoal     int
}

func makeHeapDetailsPlot(am *allMetrics) *heapDetails {
	idxobj, ok1 := am.idxs["/memory/classes/heap/objects:bytes"]
	idxunused, ok2 := am.idxs["/memory/classes/heap/unused:bytes"]
	idxfree, ok3 := am.idxs["/memory/classes/heap/free:bytes"]
	idxreleased, ok4 := am.idxs["/memory/classes/heap/released:bytes"]
	idxstacks, ok5 := am.idxs["/memory/classes/heap/stacks:bytes"]
	idxgoal, ok6 := am.idxs["/gc/heap/goal:bytes"]

	return &heapDetails{
		enabled:     ok1 && ok2 && ok3 && ok4 && ok5 && ok6,
		idxobj:      idxobj,
		idxunused:   idxunused,
		idxfree:     idxfree,
		idxreleased: idxreleased,
		idxstacks:   idxstacks,
		idxgoal:     idxgoal,
	}
}

func (p *heapDetails) name() string    { return "heap-details" }
func (p *heapDetails) isEnabled() bool { return p.enabled }

func (p *heapDetails) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:   p.name(),
		Title:  "Heap (details)",
		Type:   "scatter",
		Events: "lastgc",
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
	}
	s.Layout.Yaxis.TickSuffix = "B"
	s.Layout.Yaxis.Title = "bytes"
	return s
}

func (p *heapDetails) values(samples []metrics.Sample) interface{} {
	heapObjects := samples[p.idxobj].Value.Uint64()
	heapUnused := samples[p.idxunused].Value.Uint64()

	heapInUse := heapObjects + heapUnused
	heapFree := samples[p.idxfree].Value.Uint64()
	heapReleased := samples[p.idxreleased].Value.Uint64()

	heapIdle := heapReleased + heapFree
	heapSys := heapInUse + heapIdle
	heapStacks := samples[p.idxstacks].Value.Uint64()
	nextGC := samples[p.idxgoal].Value.Uint64()

	return []uint64{
		heapSys,
		heapObjects,
		heapStacks,
		nextGC,
	}
}

/*
 * live objects
 */

type liveObjects struct {
	enabled bool

	idxobjects int
}

func makeLiveObjectsPlot(am *allMetrics) *liveObjects {
	idxobjects, ok := am.idxs["/gc/heap/objects:objects"]

	return &liveObjects{
		enabled:    ok,
		idxobjects: idxobjects,
	}
}

func (p *liveObjects) name() string    { return "live objects" }
func (p *liveObjects) isEnabled() bool { return p.enabled }

func (p *liveObjects) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:   p.name(),
		Title:  "Live Objects in Heap",
		Type:   "bar",
		Events: "lastgc",
		Subplots: []Subplot{
			{
				Name:    "live objects",
				Unitfmt: "%{y:.4s}",
				Color:   RGBString(255, 195, 128),
			},
		},
	}
	s.Layout.Yaxis.Title = "objects"
	return s
}

func (p *liveObjects) values(samples []metrics.Sample) interface{} {
	gcHeapObjects := samples[p.idxobjects].Value.Uint64()
	return []uint64{
		gcHeapObjects,
	}
}

/*
 * live bytes
 */

type liveBytes struct {
	enabled bool

	idxallocs int
	idxfrees  int
}

func makeLiveBytesPlot(am *allMetrics) *liveBytes {
	idxallocs, ok1 := am.idxs["/gc/heap/allocs:bytes"]
	idxfrees, ok2 := am.idxs["/gc/heap/frees:bytes"]

	return &liveBytes{
		enabled:   ok1 && ok2,
		idxallocs: idxallocs,
		idxfrees:  idxfrees,
	}
}

func (p *liveBytes) name() string    { return "live bytes" }
func (p *liveBytes) isEnabled() bool { return p.enabled }

func (p *liveBytes) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:   p.name(),
		Title:  "Live Bytes in Heap",
		Type:   "bar",
		Events: "lastgc",
		Subplots: []Subplot{
			{
				Name:    "live bytes",
				Unitfmt: "%{y:.4s}B",
				Color:   RGBString(135, 182, 218),
			},
		},
	}
	s.Layout.Yaxis.Title = "bytes"
	return s
}

func (p *liveBytes) values(samples []metrics.Sample) interface{} {
	allocBytes := samples[p.idxallocs].Value.Uint64()
	freedBytes := samples[p.idxfrees].Value.Uint64()
	return []uint64{
		allocBytes - freedBytes,
	}
}

/*
 * mspan mcache
 */

type mspanMcache struct {
	enabled bool

	idxmspanInuse  int
	idxmspanFree   int
	idxmcacheInuse int
	idxmcacheFree  int
}

func makeMSpanMCachePlot(am *allMetrics) *mspanMcache {
	idxmspanInuse, ok1 := am.idxs["/memory/classes/metadata/mspan/inuse:bytes"]
	idxmspanFree, ok2 := am.idxs["/memory/classes/metadata/mspan/free:bytes"]
	idxmcacheInuse, ok3 := am.idxs["/memory/classes/metadata/mcache/inuse:bytes"]
	idxmcacheFree, ok4 := am.idxs["/memory/classes/metadata/mcache/free:bytes"]

	return &mspanMcache{
		enabled:        ok1 && ok2 && ok3 && ok4,
		idxmspanInuse:  idxmspanInuse,
		idxmspanFree:   idxmspanFree,
		idxmcacheInuse: idxmcacheInuse,
		idxmcacheFree:  idxmcacheFree,
	}
}

func (p *mspanMcache) name() string    { return "mspan-mcache" }
func (p *mspanMcache) isEnabled() bool { return p.enabled }

func (p *mspanMcache) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:   p.name(),
		Title:  "MSpan/MCache",
		Type:   "scatter",
		Events: "lastgc",
		Subplots: []Subplot{
			{
				Name:    "mspan in-use",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "mspan sys",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "mcache in-use",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "mcache sys",
				Unitfmt: "%{y:.4s}B",
			},
		},
	}
	s.Layout.Yaxis.Title = "objects"
	s.Layout.Yaxis.TickSuffix = "B"
	return s
}

func (p *mspanMcache) values(samples []metrics.Sample) interface{} {
	mspanInUse := samples[p.idxmspanInuse].Value.Uint64()
	mspanSys := samples[p.idxmspanFree].Value.Uint64()
	mcacheInUse := samples[p.idxmcacheInuse].Value.Uint64()
	mcacheSys := samples[p.idxmcacheFree].Value.Uint64()
	return []uint64{
		mspanInUse,
		mspanSys,
		mcacheInUse,
		mcacheSys,
	}
}

/*
 * goroutines
 */

type goroutines struct {
	enabled bool

	idxgs int
}

func makeGoroutinesPlot(am *allMetrics) *goroutines {
	idxgs, ok := am.idxs["/sched/goroutines:goroutines"]

	return &goroutines{
		enabled: ok,
		idxgs:   idxgs,
	}
}

func (p *goroutines) name() string    { return "goroutines" }
func (p *goroutines) isEnabled() bool { return p.enabled }

func (p *goroutines) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:   p.name(),
		Title:  "Goroutines",
		Type:   "scatter",
		Events: "lastgc",
		Subplots: []Subplot{
			{
				Name:    "goroutines",
				Unitfmt: "%{y}",
			},
		},
	}

	s.Layout.Yaxis.Title = "goroutines"
	return s
}

func (p *goroutines) values(samples []metrics.Sample) interface{} {
	return []uint64{samples[p.idxgs].Value.Uint64()}
}

/*
 * size classes
 */

type sizeClasses struct {
	enabled bool

	idxallocs int
	idxfrees  int
}

func makeSizeClassesPlot(am *allMetrics) *sizeClasses {
	idxallocs, ok1 := am.idxs["/gc/heap/allocs-by-size:bytes"]
	idxfrees, ok2 := am.idxs["/gc/heap/frees-by-size:bytes"]

	return &sizeClasses{
		enabled:   ok1 && ok2,
		idxallocs: idxallocs,
		idxfrees:  idxfrees,
	}
}

func (p *sizeClasses) name() string    { return "sizeclasses" }
func (p *sizeClasses) isEnabled() bool { return p.enabled }

func (p *sizeClasses) layout(samples []metrics.Sample) interface{} {
	// Perform a sanity check on the number of buckets on the 'allocs' and
	// 'frees' size classes histograms. Statsviz plots a single histogram based
	// on those 2 so we want them to have the same number of buckets, which
	// should be true.
	allocsBySize := samples[p.idxallocs].Value.Float64Histogram()
	freesBySize := samples[p.idxfrees].Value.Float64Histogram()
	if len(allocsBySize.Buckets) != len(freesBySize.Buckets) {
		panic("different number of buckets in allocs and frees size classes histograms!")
	}

	// No downsampling for the size classes histogram (factor=1) but we still
	// need to adapt boundaries for plotly heatmaps.
	buckets := downsampleBuckets(allocsBySize, 1)

	h := Heatmap{
		Name:       p.name(),
		Title:      "Size Classes",
		Type:       "heatmap",
		UpdateFreq: 5,
		Colorscale: BlueShades,
		Buckets:    floatseq(len(buckets)),
		CustomData: buckets,
		Hover: HeapmapHover{
			YName: "size class",
			YUnit: "bytes",
			ZName: "objects",
		},
	}
	h.Layout.Yaxis.Title = "size class"
	return h
}

func (p *sizeClasses) values(samples []metrics.Sample) interface{} {
	allocsBySize := samples[p.idxallocs].Value.Float64Histogram()
	freesBySize := samples[p.idxfrees].Value.Float64Histogram()

	sizeClasses := make([]uint64, len(allocsBySize.Counts))
	for i := 0; i < len(sizeClasses); i++ {
		sizeClasses[i] = allocsBySize.Counts[i] - freesBySize.Counts[i]
	}
	return sizeClasses
}

/*
 * gc pauses
 */

type gcpauses struct {
	enabled    bool
	histfactor int

	idxgcpauses int
}

func makeGCPausesPlot(am *allMetrics) *gcpauses {
	idxgcpauses, ok := am.idxs["/gc/pauses:seconds"]

	return &gcpauses{
		enabled:     ok,
		idxgcpauses: idxgcpauses,
	}
}

func (p *gcpauses) name() string    { return "gcpauses" }
func (p *gcpauses) isEnabled() bool { return p.enabled }

func (p *gcpauses) layout(samples []metrics.Sample) interface{} {
	gcpauses := samples[p.idxgcpauses].Value.Float64Histogram()
	p.histfactor = downsampleFactor(len(gcpauses.Buckets), maxBuckets)
	buckets := downsampleBuckets(gcpauses, p.histfactor)

	h := Heatmap{
		Name:       p.name(),
		Title:      "Stop-the-world pause latencies",
		Type:       "heatmap",
		UpdateFreq: 5,
		Colorscale: PinkShades,
		Buckets:    floatseq(len(buckets)),
		CustomData: buckets,
		Hover: HeapmapHover{
			YName: "pause duration",
			YUnit: "duration",
			ZName: "pauses",
		},
	}
	h.Layout.Yaxis.Title = "pause duration"
	return h
}

func (p *gcpauses) values(samples []metrics.Sample) interface{} {
	gcpauses := samples[p.idxgcpauses].Value.Float64Histogram()
	return downsampleCounts(gcpauses, p.histfactor)
}

/*
 * scheduler latencies
 */

type schedlat struct {
	enabled    bool
	histfactor int

	idxschedlat int
}

func makeSchedLatPlot(am *allMetrics) *schedlat {
	idxschedlat, ok := am.idxs["/sched/latencies:seconds"]

	return &schedlat{
		enabled:     ok,
		idxschedlat: idxschedlat,
	}
}

func (p *schedlat) name() string    { return "sched-latencies" }
func (p *schedlat) isEnabled() bool { return p.enabled }

func (p *schedlat) layout(samples []metrics.Sample) interface{} {
	schedlat := samples[p.idxschedlat].Value.Float64Histogram()
	p.histfactor = downsampleFactor(len(schedlat.Buckets), maxBuckets)
	buckets := downsampleBuckets(schedlat, p.histfactor)

	h := Heatmap{
		Name:       p.name(),
		Title:      "Time in scheduler before a goroutine runs",
		Type:       "heatmap",
		UpdateFreq: 5,
		Colorscale: GreenShades,
		Buckets:    floatseq(len(buckets)),
		CustomData: buckets,
		Hover: HeapmapHover{
			YName: "duration",
			YUnit: "duration",
			ZName: "goroutines",
		},
	}
	h.Layout.Yaxis.Title = "duration"
	return h
}

func (p *schedlat) values(samples []metrics.Sample) interface{} {
	schedlat := samples[p.idxschedlat].Value.Float64Histogram()
	return downsampleCounts(schedlat, p.histfactor)
}
