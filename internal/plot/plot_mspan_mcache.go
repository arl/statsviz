package plot

import (
	"runtime/metrics"
	"time"
)

var _ = register(description{
	metrics: []string{
		"/memory/classes/metadata/mspan/inuse:bytes",
		"/memory/classes/metadata/mspan/free:bytes",
		"/memory/classes/metadata/mcache/inuse:bytes",
		"/memory/classes/metadata/mcache/free:bytes",
	},
	getvalues: func() getvalues {
		return func(_ time.Time, samples []metrics.Sample) any {
			mspanInUse := samples[idx_memory_classes_metadata_mspan_inuse_bytes].Value.Uint64()
			mspanFree := samples[idx_memory_classes_metadata_mspan_free_bytes].Value.Uint64()
			mcacheInUse := samples[idx_memory_classes_metadata_mcache_inuse_bytes].Value.Uint64()
			mcacheFree := samples[idx_memory_classes_metadata_mcache_free_bytes].Value.Uint64()

			return []uint64{
				mspanInUse,
				mspanFree,
				mcacheInUse,
				mcacheFree,
			}
		}
	},
	layout: Scatter{
		Name:   "mspan-mcache",
		Tags:   []tag{tagGC},
		Title:  "MSpan/MCache",
		Type:   "scatter",
		Events: "lastgc",
		Layout: ScatterLayout{
			Yaxis: ScatterYAxis{
				Title:      "bytes",
				TickSuffix: "B",
			},
		},
		Subplots: []Subplot{
			{Name: "mspan in-use", Unitfmt: "%{y:.4s}B"},
			{Name: "mspan free", Unitfmt: "%{y:.4s}B"},
			{Name: "mcache in-use", Unitfmt: "%{y:.4s}B"},
			{Name: "mcache free", Unitfmt: "%{y:.4s}B"},
		},
		InfoText: `<i>mspan in-use</i> is <b>/memory/classes/metadata/mspan/inuse:bytes</b>. Memory that is occupied by runtime mspan structures that are currently being used.
<i>mspan free</i> is <b>/memory/classes/metadata/mspan/free:bytes</b>. Memory that is reserved for runtime mspan structures, but not in-use.
<i>mcache in-use</i> is <b>/memory/classes/metadata/mcache/inuse:bytes</b>. Memory that is occupied by runtime mcache structures that are currently being used.
<i>mcache free</i> is <b>/memory/classes/metadata/mcache/free:bytes</b>. Memory that is reserved for runtime mcache structures, but not in-use.`,
	},
})
