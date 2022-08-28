package plot

import (
	"encoding/json"
	"fmt"
	"io"
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

func (pl *List) initMetrics() {
	pl.idxs = make(map[string]int)
	pl.descs = metrics.All()
	pl.samples = make([]metrics.Sample, len(pl.descs))
	for i := range pl.samples {
		pl.samples[i].Name = pl.descs[i].Name
		pl.idxs[pl.samples[i].Name] = i
	}
}

func (pl *List) Config() *Config {
	pl.once.Do(pl.config)
	return pl.cfg
}

func (pl *List) config() {
	pl.initMetrics()

	pl.plots = append(pl.plots, makeHeapGlobalPlot(pl.idxs))
	pl.plots = append(pl.plots, makeHeapDetailsPlot(pl.idxs))
	pl.plots = append(pl.plots, makeLiveObjectsPlot(pl.idxs))
	pl.plots = append(pl.plots, makeLiveBytesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeMSpanMCachePlot(pl.idxs))
	pl.plots = append(pl.plots, makeGoroutinesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeSizeClassesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeGCPausesPlot(pl.idxs))
	pl.plots = append(pl.plots, makeSchedLatPlot(pl.idxs))
	pl.plots = append(pl.plots, makeGCStackSize(pl.idxs))
	pl.plots = append(pl.plots, makeSchedEvents(pl.idxs))
	pl.plots = append(pl.plots, makeCGOPlot(pl.idxs))

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
