package plot

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"runtime/metrics"
	"sync"
)

// allMetrics contains the descriptions and samples for all suported metrics (as
// per metrics.All()).
type allMetrics struct {
	idxs    map[string]int // metric name -> index in descs and samples
	descs   []metrics.Description
	samples []metrics.Sample
}

func (am *allMetrics) init() {
	am.idxs = make(map[string]int)
	am.descs = metrics.All()
	am.samples = make([]metrics.Sample, len(am.descs))
	for i := range am.samples {
		am.samples[i].Name = am.descs[i].Name
		am.idxs[am.samples[i].Name] = i
	}
}

type plot interface {
	name() string
	isEnabled() bool
	layout([]metrics.Sample) interface{}
	values([]metrics.Sample) interface{}
}

var All List

type List struct {
	plots []plot

	once sync.Once
	cfg  *Config

	mu sync.Mutex
	am allMetrics
}

func (pl *List) Config() *Config {
	pl.once.Do(func() {
		pl.am.init()

		pl.plots = append(pl.plots, makeHeapGlobalPlot(&pl.am))
		pl.plots = append(pl.plots, makeHeapDetailsPlot(&pl.am))
		pl.plots = append(pl.plots, makeLiveObjectsPlot(&pl.am))
		pl.plots = append(pl.plots, makeLiveBytesPlot(&pl.am))
		pl.plots = append(pl.plots, makeMSpanMCachePlot(&pl.am))
		pl.plots = append(pl.plots, makeGoroutinesPlot(&pl.am))
		pl.plots = append(pl.plots, makeSizeClassesPlot(&pl.am))
		pl.plots = append(pl.plots, makeGCPausesPlot(&pl.am))
		pl.plots = append(pl.plots, makeSchedLatPlot(&pl.am))

		metrics.Read(pl.am.samples)

		var layouts []interface{}
		for _, p := range pl.plots {
			if p.isEnabled() {
				layouts = append(layouts, p.layout(pl.am.samples))
			}
		}

		pl.cfg = &Config{
			Events: []string{"lastgc"},
			Series: layouts,
		}
	})
	return pl.cfg
}

func (pl *List) WriteValues(w io.Writer) error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	metrics.Read(pl.am.samples)

	m := make(map[string]interface{})
	for _, p := range pl.plots {
		if p.isEnabled() {
			m[p.name()] = p.values(pl.am.samples)
		}
	}

	// lastgc time series is used as source to represent garbage collection
	// timestamps as vertical bars on certain plots.
	gcStats := debug.GCStats{}
	debug.ReadGCStats(&gcStats)
	// In javascript, timestamps are in ms.
	lastgc := gcStats.LastGC.UnixMilli()
	m["lastgc"] = []int64{lastgc}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		return fmt.Errorf("failed to write/convert metrics values to json: %v", err)
	}
	return nil
}
