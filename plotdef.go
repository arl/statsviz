package statsviz

import (
	"fmt"
	"image/color"
	"runtime"
)

// Those are the JSON structures and should be placed under
// internal so as to not be part of the public API.

type (
	PlotsDefinition struct {
		Events []string      `json:"events"`
		Series []interface{} `json:"series"`
	}

	ScatterPlotLayout struct {
		Yaxis ScatterPlotLayoutYAxis `json:"yaxis"`
	}

	ScatterPlotLayoutYAxis struct {
		Title      string `json:"title"`
		TickSuffix string `json:"ticksuffix"`
		TickFormat string `json:"tickformat"`
	}

	ScatterPlotSubplot struct {
		Name    string `json:"name"`
		Hover   string `json:"hover"`
		Unitfmt string `json:"unitfmt"`
	}

	ScatterPlot struct {
		Name       string               `json:"name"`
		Title      string               `json:"title"`
		Type       string               `json:"type"`
		UpdateFreq int                  `json:"updateFreq"`
		HorzEvents string               `json:"horzEvents"`
		Layout     ScatterPlotLayout    `json:"layout"`
		Subplots   []ScatterPlotSubplot `json:"subplots"`
	}

	HeatmapPlot struct {
		Name       string            `json:"name"`
		Title      string            `json:"title"`
		Type       string            `json:"type"`
		UpdateFreq int               `json:"updateFreq"`
		HorzEvents string            `json:"horzEvents"`
		Layout     HeatmapPlotLayout `json:"layout"`
		Heatmap    Heatmap           `json:"heatmap"`
	}

	HeatmapPlotLayout struct {
		Yaxis HeatmapPlotLayoutYAxis `json:"yaxis"`
	}

	HeatmapPlotLayoutYAxis struct {
		Title string `json:"title"`
	}

	Heatmap struct {
		Hover      string    `json:"hover"`
		Colorscale []Color   `json:"colorscale"`
		Buckets    []float64 `json:"buckets"`
	}

	Color struct {
		Max   float64
		Color color.RGBA
	}
)

func (c Color) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf(`[%f,"rgb(%d,%d,%d,%f)"]`,
		c.Max, c.Color.R, c.Color.G, c.Color.B, float64(c.Color.A)/255)
	return []byte(str), nil
}

// TODO(arl) rename
type plotsDefBuilder struct {
	series []interface{}
}

// TODO(arl) should probably transofrm this function so that
// it just takes fields and not the whole ScatterPlot
func (b *plotsDefBuilder) addScatter(scatter ScatterPlot) {
	b.series = append(b.series, scatter)
}

func (b *plotsDefBuilder) addHeatmap(heatmap HeatmapPlot) {
	b.series = append(b.series, heatmap)
}

func (b *plotsDefBuilder) close() PlotsDefinition {
	def := PlotsDefinition{
		Series: []interface{}{},
		Events: []string{},
	}
	horzEvents := make(map[string]struct{})

	for i := range b.series {
		def.Series = append(def.Series, b.series[i])
		switch val := b.series[i].(type) {
		case ScatterPlot:
			horzEvents[val.HorzEvents] = struct{}{}
			val.Type = "scatter"
			b.series[i] = val
		case HeatmapPlot:
			horzEvents[val.HorzEvents] = struct{}{}
			val.Type = "heatmap"
			b.series[i] = val
		}
	}

	for e := range horzEvents {
		if e != "" {
			def.Events = append(def.Events, e)
		}
	}

	return def
}

// CONTINUER ICI

type Axis struct {
	Title string
	Unit  Unit
}

type Unit struct {
	TickSuffix string
	UnitFmt    string
}

var (
	Bytes = Unit{TickSuffix: "B", UnitFmt: "%{y:.4s}B"}
)

var plotsDef = PlotsDefinition{
	Events: []string{"lastgc"},
	Series: []interface{}{
		ScatterPlot{
			Name:       "heap",
			Title:      "Heap",
			Type:       "scatter",
			UpdateFreq: 0,
			HorzEvents: "lastgc",
			Layout: ScatterPlotLayout{
				Yaxis: ScatterPlotLayoutYAxis{
					Title:      "bytes",
					TickSuffix: "B",
				},
			},
			Subplots: []ScatterPlotSubplot{
				{
					Name:    "heap alloc",
					Hover:   "heap alloc",
					Unitfmt: "%{y:.4s}B",
				},
				{
					Name:    "heap sys",
					Hover:   "heap sys",
					Unitfmt: "%{y:.4s}B",
				},
				{
					Name:    "heap idle",
					Hover:   "heap idle",
					Unitfmt: "%{y:.4s}B",
				},
				{
					Name:    "heap in-use",
					Hover:   "heap in-use",
					Unitfmt: "%{y:.4s}B",
				},
				{
					Name:    "heap next gc",
					Hover:   "heap next gc",
					Unitfmt: "%{y:.4s}B",
				},
			},
		},
		ScatterPlot{
			Name:       "objects",
			Title:      "Objects",
			Type:       "scatter",
			UpdateFreq: 0,
			HorzEvents: "lastgc",
			Layout: ScatterPlotLayout{
				Yaxis: ScatterPlotLayoutYAxis{
					Title: "objects",
				},
			},
			Subplots: []ScatterPlotSubplot{
				{
					Name:    "live",
					Hover:   "live objects",
					Unitfmt: "%{y}",
				},
				{
					Name:    "lookups",
					Hover:   "pointer lookups",
					Unitfmt: "%{y}",
				},
				{
					Name:    "heap",
					Hover:   "heap objects",
					Unitfmt: "%{y}",
				},
			},
		},
		ScatterPlot{
			Name:       "mspan-mcache",
			Title:      "MSpan/MCache",
			Type:       "scatter",
			UpdateFreq: 0,
			HorzEvents: "lastgc",
			Layout: ScatterPlotLayout{
				Yaxis: ScatterPlotLayoutYAxis{
					Title:      "bytes",
					TickSuffix: "B",
				},
			},
			Subplots: []ScatterPlotSubplot{
				{
					Name:    "mspan in-use",
					Hover:   "mspan in-use",
					Unitfmt: "%{y:.4s}B",
				},
				{
					Name:    "mspan sys",
					Hover:   "mspan sys",
					Unitfmt: "%{y:.4s}B",
				},
				{
					Name:    "mcache in-use",
					Hover:   "mcache in-use",
					Unitfmt: "%{y:.4s}B",
				},
				{
					Name:    "mcache sys",
					Hover:   "mcache sys",
					Unitfmt: "%{y:.4s}B",
				},
			},
		},
		ScatterPlot{
			Name:       "goroutines",
			Title:      "Goroutines",
			Type:       "scatter",
			UpdateFreq: 0,
			HorzEvents: "lastgc",
			Layout: ScatterPlotLayout{
				Yaxis: ScatterPlotLayoutYAxis{
					Title: "goroutines",
				},
			},
			Subplots: []ScatterPlotSubplot{
				{
					Name:    "goroutines",
					Unitfmt: "%{y}",
				},
			},
		},
		/* TODO: Heatmap */

		ScatterPlot{
			Name:       "gcfraction",
			Title:      "GC/CPU fraction",
			Type:       "scatter",
			UpdateFreq: 0,
			HorzEvents: "lastgc",
			Layout: ScatterPlotLayout{
				Yaxis: ScatterPlotLayoutYAxis{
					Title:      "gc/cpu (%)",
					TickFormat: ",.5%",
				},
			},
			Subplots: []ScatterPlotSubplot{
				{
					Name:    "gc/cpu",
					Hover:   "gc/cpu fraction",
					Unitfmt: "%{y:,.4%}",
				},
			},
		},
	},
}

func plotsValues() map[string]interface{} {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	numgs := runtime.NumGoroutine()

	m := make(map[string]interface{})
	m["heap"] = []uint64{stats.HeapAlloc, stats.HeapSys, stats.HeapIdle, stats.HeapInuse, stats.NextGC}
	m["objects"] = []uint64{stats.Alloc - stats.Frees, stats.Lookups, stats.HeapObjects}
	m["mspan-mcache"] = []uint64{stats.MSpanInuse, stats.MSpanSys, stats.MCacheInuse, stats.MCacheSys}
	m["goroutines"] = []int{numgs}
	m["gcfraction"] = []float64{stats.GCCPUFraction}
	m["lastgc"] = []uint64{stats.LastGC / 1_000_000} // Javascript datetime is in ms
	return m
}
