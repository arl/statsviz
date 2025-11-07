package plot

import "runtime/metrics"

var _ = register(description{
	name: "mspan-mcache",
	tags: []tag{tagGC},
	metrics: []string{
		"/memory/classes/metadata/mspan/inuse:bytes",
		"/memory/classes/metadata/mspan/free:bytes",
		"/memory/classes/metadata/mcache/inuse:bytes",
		"/memory/classes/metadata/mcache/free:bytes",
	},
	layout: Scatter{
		Name:   "TODO(set later)",
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
		InfoText: `
<i>Mspan in-use</i> is <b>/memory/classes/metadata/mspan/inuse</b>, the memory that is occupied by runtime mspan structures that are currently being used.
<i>Mspan free</i> is <b>/memory/classes/metadata/mspan/free</b>, the memory that is reserved for runtime mspan structures, but not in-use.
<i>Mcache in-use</i> is <b>/memory/classes/metadata/mcache/inuse</b>, the memory that is occupied by runtime mcache structures that are currently being used.
<i>Mcache free</i> is <b>/memory/classes/metadata/mcache/free</b>, the memory that is reserved for runtime mcache structures, but not in-use.
`,
	},
	make: func(indices ...int) metricsGetter {
		return &mspanMcache{
			idxmspanInuse:  indices[0],
			idxmspanFree:   indices[1],
			idxmcacheInuse: indices[2],
			idxmcacheFree:  indices[3],
		}
	},
})

type mspanMcache struct {
	enabled bool

	idxmspanInuse  int
	idxmspanFree   int
	idxmcacheInuse int
	idxmcacheFree  int
}

func (p *mspanMcache) values(samples []metrics.Sample) any {
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
