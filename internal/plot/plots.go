package plot

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"runtime/debug"
	"runtime/metrics"
	"sync"
	"time"
)

type plot interface {
	name() string
	isEnabled() bool
	layout([]metrics.Sample) interface{}
	values([]metrics.Sample) interface{}
}

// List holds all the plots that statsviz knows about. Some plots might be
// disabled, if they rely on metrics that are unknown to the current Go version.
type List struct {
	plots []plot

	once sync.Once // ensure Config is called once
	cfg  *Config

	idxs  map[string]int // map metrics name to idx in samples and descs
	descs []metrics.Description

	mu      sync.Mutex // protects samples in case of concurrent calls to WriteValues
	samples []metrics.Sample
}

func NewList() *List {
	descs := metrics.All()
	pl := &List{
		idxs:    make(map[string]int),
		descs:   descs,
		samples: make([]metrics.Sample, len(descs)),
	}

	for i := range pl.samples {
		pl.samples[i].Name = pl.descs[i].Name
		pl.idxs[pl.samples[i].Name] = i
	}

	pl.addRuntimeMetrics()
	return pl
}

func (pl *List) Config() *Config {
	pl.once.Do(pl.genConfig)
	return pl.cfg
}

func (pl *List) genConfig() {
	metrics.Read(pl.samples)

	var layouts []interface{}
	for _, p := range pl.plots {
		if p.isEnabled() {
			layouts = append(layouts, p.layout(pl.samples))
		}
	}

	pl.cfg = &Config{
		Events: []string{"lastgc"},
		Series: layouts,
	}
}

func (pl *List) addRuntimeMetrics() {
	pl.plots = append(pl.plots, makeHeapGlobalPlot(pl.idxs))
	pl.plots = append(pl.plots, makeHeapDetailsPlot(pl.idxs))
	pl.plots = append(pl.plots, makeLiveObjectsPlot(pl.idxs))
	pl.plots = append(pl.plots, makeLiveBytesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeMSpanMCachePlot(pl.idxs))
	pl.plots = append(pl.plots, makeGoroutinesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeSizeClassesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeGCPausesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeRunnableTime(pl.idxs))
	pl.plots = append(pl.plots, makeGCStackSize(pl.idxs))
	pl.plots = append(pl.plots, makeSchedEvents(pl.idxs))
	pl.plots = append(pl.plots, makeCGOPlot(pl.idxs))
}

// WriteValues writes into w a JSON object containing the data points for all
// plots at the current instant.
func (pl *List) WriteValues(w io.Writer) error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	metrics.Read(pl.samples)

	// lastgc time series is used as source to represent garbage collection
	// timestamps as vertical bars on certain plots.
	gcStats := debug.GCStats{}
	debug.ReadGCStats(&gcStats)

	m := make(map[string]interface{})
	for _, p := range pl.plots {
		if p.isEnabled() {
			m[p.name()] = p.values(pl.samples)
		}
	}
	// In javascript, timestamps are in ms.
	m["lastgc"] = []int64{gcStats.LastGC.UnixMilli()}
	m["timestamp"] = time.Now().UnixMilli()

	if err := json.NewEncoder(w).Encode(m); err != nil {
		return fmt.Errorf("failed to write/convert metrics values to json: %v", err)
	}
	return nil
}

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

func makeHeapGlobalPlot(idxs map[string]int) *heapGlobal {
	idxobj, ok1 := idxs["/memory/classes/heap/objects:bytes"]
	idxunused, ok2 := idxs["/memory/classes/heap/unused:bytes"]
	idxfree, ok3 := idxs["/memory/classes/heap/free:bytes"]
	idxreleased, ok4 := idxs["/memory/classes/heap/released:bytes"]

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
		InfoText: `<i>Heap in use</i> is <b>/memory/classes/heap/objects + /memory/classes/heap/unused</b>. It amounts to the memory occupied by live objects and dead objects that are not yet marked free by the GC, plus some memory reserved for heap objects.
<i>Heap free</i> is <b>/memory/classes/heap/free</b>, that is free memory that could be returned to the OS, but has not been.
<i>Heap released</i> is <b>/memory/classes/heap/free</b>, memory that is free memory that has been returned to the OS.`,
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

func makeHeapDetailsPlot(idxs map[string]int) *heapDetails {
	idxobj, ok1 := idxs["/memory/classes/heap/objects:bytes"]
	idxunused, ok2 := idxs["/memory/classes/heap/unused:bytes"]
	idxfree, ok3 := idxs["/memory/classes/heap/free:bytes"]
	idxreleased, ok4 := idxs["/memory/classes/heap/released:bytes"]
	idxstacks, ok5 := idxs["/memory/classes/heap/stacks:bytes"]
	idxgoal, ok6 := idxs["/gc/heap/goal:bytes"]

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
		InfoText: `<i>Heap</i> sys is <b>/memory/classes/heap/objects + /memory/classes/heap/unused + /memory/classes/heap/released + /memory/classes/heap/free</b>. It's an estimate of all the heap memory obtained form the OS.
<i>Heap objects</i> is <b>/memory/classes/heap/objects</b>, the memory occupied by live objects and dead objects that have not yet been marked free by the GC.
<i>Heap stacks</i> is <b>/memory/classes/heap/stacks</b>, the memory used for stack space.
<i>Heap goal</i> is <b>gc/heap/goal</b>, the heap size target for the end of the GC cycle.`,
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

func makeLiveObjectsPlot(idxs map[string]int) *liveObjects {
	idxobjects, ok := idxs["/gc/heap/objects:objects"]

	return &liveObjects{
		enabled:    ok,
		idxobjects: idxobjects,
	}
}

func (p *liveObjects) name() string    { return "live-objects" }
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
		InfoText: `<i>Live objects</i> is <b>/gc/heap/objects</b>. It's the number of objects, live or unswept, occupying heap memory.`,
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

func makeLiveBytesPlot(idxs map[string]int) *liveBytes {
	idxallocs, ok1 := idxs["/gc/heap/allocs:bytes"]
	idxfrees, ok2 := idxs["/gc/heap/frees:bytes"]

	return &liveBytes{
		enabled:   ok1 && ok2,
		idxallocs: idxallocs,
		idxfrees:  idxfrees,
	}
}

func (p *liveBytes) name() string    { return "live-bytes" }
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
		InfoText: `<i>Live bytes</i> is <b>/gc/heap/allocs - /gc/heap/frees</b>. It's the number of bytes currently allocated (and not yet GC'ec) to the heap by the application.`,
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

func makeMSpanMCachePlot(idxs map[string]int) *mspanMcache {
	idxmspanInuse, ok1 := idxs["/memory/classes/metadata/mspan/inuse:bytes"]
	idxmspanFree, ok2 := idxs["/memory/classes/metadata/mspan/free:bytes"]
	idxmcacheInuse, ok3 := idxs["/memory/classes/metadata/mcache/inuse:bytes"]
	idxmcacheFree, ok4 := idxs["/memory/classes/metadata/mcache/free:bytes"]

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
				Name:    "mspan free",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "mcache in-use",
				Unitfmt: "%{y:.4s}B",
			},
			{
				Name:    "mcache free",
				Unitfmt: "%{y:.4s}B",
			},
		},
		InfoText: `<i>Mspan in-use</i> is <b>/memory/classes/metadata/mspan/inuse</b>, the memory that is occupied by runtime mspan structures that are currently being used.
<i>Mspan free</i> is <b>/memory/classes/metadata/mspan/free</b>, the memory that is reserved for runtime mspan structures, but not in-use.
<i>Mcache in-use</i> is <b>/memory/classes/metadata/mcache/inuse</b>, the memory that is occupied by runtime mcache structures that are currently being used.
<i>Mcache free</i> is <b>/memory/classes/metadata/mcache/free</b>, the memory that is reserved for runtime mcache structures, but not in-use.
`,
	}
	s.Layout.Yaxis.Title = "bytes"
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

func makeGoroutinesPlot(idxs map[string]int) *goroutines {
	idxgs, ok := idxs["/sched/goroutines:goroutines"]

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
		InfoText: "<i>Goroutines</i> is <b>/sched/goroutines</b>, the count of live goroutines.",
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
	enabled     bool
	sizeClasses []uint64

	idxallocs int
	idxfrees  int
}

func makeSizeClassesPlot(idxs map[string]int) *sizeClasses {
	idxallocs, ok1 := idxs["/gc/heap/allocs-by-size:bytes"]
	idxfrees, ok2 := idxs["/gc/heap/frees-by-size:bytes"]

	return &sizeClasses{
		enabled:   ok1 && ok2,
		idxallocs: idxallocs,
		idxfrees:  idxfrees,
	}
}

func (p *sizeClasses) name() string    { return "size-classes" }
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

	// Pre-allocate here so we never do it in values.
	p.sizeClasses = make([]uint64, len(allocsBySize.Counts))

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
		InfoText: `This heatmap shows the distribution of size classes, using <b>/gc/heap/allocs-by-size</b> and <b>/gc/heap/frees-by-size</b>.`,
	}
	h.Layout.Yaxis.Title = "size class"
	h.Layout.Yaxis.TickMode = "array"
	h.Layout.Yaxis.TickVals = []float64{1, 9, 17, 25, 31, 37, 43, 50, 58, 66}
	h.Layout.Yaxis.TickText = []float64{1 << 4, 1 << 7, 1 << 8, 1 << 9, 1 << 10, 1 << 11, 1 << 12, 1 << 13, 1 << 14, 1 << 15}
	return h
}

func (p *sizeClasses) values(samples []metrics.Sample) interface{} {
	allocsBySize := samples[p.idxallocs].Value.Float64Histogram()
	freesBySize := samples[p.idxfrees].Value.Float64Histogram()

	for i := 0; i < len(p.sizeClasses); i++ {
		p.sizeClasses[i] = allocsBySize.Counts[i] - freesBySize.Counts[i]
	}
	return p.sizeClasses
}

/*
 * gc pauses
 */

type gcpauses struct {
	enabled    bool
	histfactor int
	counts     [maxBuckets]uint64

	idxgcpauses int
}

func makeGCPausesPlot(idxs map[string]int) *gcpauses {
	idxgcpauses, ok := idxs["/gc/pauses:seconds"]

	return &gcpauses{
		enabled:     ok,
		idxgcpauses: idxgcpauses,
	}
}

func (p *gcpauses) name() string    { return "gc-pauses" }
func (p *gcpauses) isEnabled() bool { return p.enabled }

func (p *gcpauses) layout(samples []metrics.Sample) interface{} {
	gcpauses := samples[p.idxgcpauses].Value.Float64Histogram()
	p.histfactor = downsampleFactor(len(gcpauses.Buckets), maxBuckets)
	buckets := downsampleBuckets(gcpauses, p.histfactor)

	h := Heatmap{
		Name:       p.name(),
		Title:      "Stop-the-world Pause Latencies",
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
		InfoText: `This heatmap shows the distribution of individual GC-related stop-the-world pause latencies, uses <b>/gc/pauses:seconds</b>,.`,
	}
	h.Layout.Yaxis.Title = "pause duration"
	h.Layout.Yaxis.TickMode = "array"
	h.Layout.Yaxis.TickVals = []float64{6, 13, 20, 26, 33, 39.5, 46, 53, 60, 66, 73, 79, 86}
	ticks := []float64{-7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5}
	for _, tick := range ticks {
		h.Layout.Yaxis.TickText = append(h.Layout.Yaxis.TickText, math.Pow(10, tick))
	}
	return h
}

func (p *gcpauses) values(samples []metrics.Sample) interface{} {
	gcpauses := samples[p.idxgcpauses].Value.Float64Histogram()
	return downsampleCounts(gcpauses, p.histfactor, p.counts[:])
}

/*
 * time spent in runnable state
 */

type runnableTime struct {
	enabled    bool
	histfactor int
	counts     [maxBuckets]uint64

	idxschedlat int
}

func makeRunnableTime(idxs map[string]int) *runnableTime {
	idxschedlat, ok := idxs["/sched/latencies:seconds"]

	return &runnableTime{
		enabled:     ok,
		idxschedlat: idxschedlat,
	}
}

func (p *runnableTime) name() string    { return "runnable-time" }
func (p *runnableTime) isEnabled() bool { return p.enabled }

func (p *runnableTime) layout(samples []metrics.Sample) interface{} {
	schedlat := samples[p.idxschedlat].Value.Float64Histogram()
	p.histfactor = downsampleFactor(len(schedlat.Buckets), maxBuckets)
	buckets := downsampleBuckets(schedlat, p.histfactor)

	h := Heatmap{
		Name:       p.name(),
		Title:      "Time Goroutines Spend in 'Runnable'",
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
		InfoText: `This heatmap shows the distribution of the time goroutines have spent in the scheduler in a runnable state before actually running, uses <b>/sched/latencies:seconds</b>.`,
	}
	h.Layout.Yaxis.Title = "duration"
	h.Layout.Yaxis.TickMode = "array"
	h.Layout.Yaxis.TickVals = []float64{6, 13, 20, 26, 33, 39.5, 46, 53, 60, 66, 73, 79, 86}
	ticks := []float64{-7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5}
	for _, tick := range ticks {
		h.Layout.Yaxis.TickText = append(h.Layout.Yaxis.TickText, math.Pow(10, tick))
	}

	return h
}

func (p *runnableTime) values(samples []metrics.Sample) interface{} {
	schedlat := samples[p.idxschedlat].Value.Float64Histogram()

	return downsampleCounts(schedlat, p.histfactor, p.counts[:])
}

/*
 * scheduling events
 */

type schedEvents struct {
	enabled bool

	idxschedlat   int
	idxGomaxprocs int
	lasttot       uint64
}

func makeSchedEvents(idxs map[string]int) *schedEvents {
	idxschedlat, ok1 := idxs["/sched/latencies:seconds"]
	idxGomaxprocs, ok2 := idxs["/sched/gomaxprocs:threads"]

	return &schedEvents{
		enabled:       ok1 && ok2,
		idxschedlat:   idxschedlat,
		idxGomaxprocs: idxGomaxprocs,
		lasttot:       math.MaxUint64,
	}
}

func (p *schedEvents) name() string    { return "sched-events" }
func (p *schedEvents) isEnabled() bool { return p.enabled }

func (p *schedEvents) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:   p.name(),
		Title:  "Goroutine Scheduling Events",
		Type:   "scatter",
		Events: "lastgc",
		Subplots: []Subplot{
			{
				Name:    "events per unit of time",
				Unitfmt: "%{y}",
			},
			{
				Name:    "events per unit of time, per P",
				Unitfmt: "%{y}",
			},
		},
		InfoText: `<i>Events per second</i> is the sum of all buckets in <b>/sched/latencies:seconds</b>, that is, it tracks the total number of goroutine scheduling events. That number is multiplied by the constant 8.
<i>Events per second per P (processor)</i> is <i>Events per second</i> divided by current <b>GOMAXPROCS</b>, from <b>/sched/gomaxprocs:threads</b>.
<b>NOTE</b>: the multiplying factor comes from internal Go runtime source code and might change from version to version.`,
	}
	s.Layout.Yaxis.Title = "events"
	return s
}

// gTrackingPeriod is currently always 8. Guard it behind build tags when that
// changes. See https://github.com/golang/go/blob/go1.18.4/src/runtime/runtime2.go#L502-L504
const currentGtrackingPeriod = 8

func (p *schedEvents) values(samples []metrics.Sample) interface{} {
	schedlat := samples[p.idxschedlat].Value.Float64Histogram()
	gomaxprocs := samples[p.idxGomaxprocs].Value.Uint64()

	total := uint64(0)
	for _, v := range schedlat.Counts {
		total += v
	}
	total *= currentGtrackingPeriod

	curtot := total - p.lasttot
	if p.lasttot == math.MaxUint64 {
		// We don't want a big spike at statsviz launch in case the process has
		// been running for some time and curtot is high.
		curtot = 0
	}
	p.lasttot = total

	ftot := float64(curtot)

	return []float64{
		ftot,
		ftot / float64(gomaxprocs),
	}
}

/*
 * cgo
 */

type cgo struct {
	enabled  bool
	idxgo2c  int
	lastgo2c uint64
}

func makeCGOPlot(idxs map[string]int) *cgo {
	idxgo2c, ok := idxs["/cgo/go-to-c-calls:calls"]

	return &cgo{
		enabled:  ok,
		idxgo2c:  idxgo2c,
		lastgo2c: math.MaxUint64,
	}
}

func (p *cgo) name() string    { return "cgo" }
func (p *cgo) isEnabled() bool { return p.enabled }

func (p *cgo) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:  p.name(),
		Title: "CGO Calls",
		Type:  "bar",
		Subplots: []Subplot{
			{
				Name:    "calls from go to c",
				Unitfmt: "%{y}",
				Color:   "red",
			},
		},
		InfoText: "Shows the count of calls made from Go to C by the current process, per unit of time. Uses <b>/cgo/go-to-c-calls:calls</b>",
	}

	s.Layout.Yaxis.Title = "calls"
	return s
}

func (p *cgo) values(samples []metrics.Sample) interface{} {
	go2c := samples[p.idxgo2c].Value.Uint64()
	curgo2c := go2c - p.lastgo2c
	if p.lastgo2c == math.MaxUint64 {
		// We don't want a big spike at statsviz launch in case the process has
		// been running for some time and curgo2c is high.
		curgo2c = 0
	}
	p.lastgo2c = go2c

	return []uint64{curgo2c}
}

/*
 * gc stack size
 */

type gcStackSize struct {
	enabled  bool
	idxstack int
}

func makeGCStackSize(idxs map[string]int) *gcStackSize {
	idxstack, ok := idxs["/gc/stack/starting-size:bytes"]

	return &gcStackSize{
		enabled:  ok,
		idxstack: idxstack,
	}
}

func (p *gcStackSize) name() string    { return "gc-stack-size" }
func (p *gcStackSize) isEnabled() bool { return p.enabled }

func (p *gcStackSize) layout(_ []metrics.Sample) interface{} {
	s := Scatter{
		Name:  p.name(),
		Title: "Starting Size of Goroutines Stacks",
		Type:  "scatter",
		Subplots: []Subplot{
			{
				Name:    "new goroutines stack size",
				Unitfmt: "%{y:.4s}B",
			},
		},
		InfoText: "Shows the stack size of new goroutines, uses <b>/gc/stack/starting-size:bytes</b>",
	}

	s.Layout.Yaxis.Title = "bytes"
	return s
}

func (p *gcStackSize) values(samples []metrics.Sample) interface{} {
	stackSize := samples[p.idxstack].Value.Uint64()
	return []uint64{stackSize}
}

/*
 * helpers
 */

func floatseq(n int) []float64 {
	seq := make([]float64, n)
	for i := 0; i < n; i++ {
		seq[i] = float64(i)
	}
	return seq
}
