package plot

import "runtime/metrics"

type stoppingPausesGC struct {
	histfactor int
	counts     [maxBuckets]uint64

	idxstoppinggc int
}

func makeStoppingPausesGC(indices ...int) metricsGetter {
	return &stoppingPausesGC{
		idxstoppinggc: indices[0],
	}
}

func (p *stoppingPausesGC) values(samples []metrics.Sample) any {
	if p.histfactor == 0 {
		stoppinggc := samples[p.idxstoppinggc].Value.Float64Histogram()
		p.histfactor = downsampleFactor(len(stoppinggc.Buckets), maxBuckets)
	}

	stoppinggc := samples[p.idxstoppinggc].Value.Float64Histogram()
	return downsampleCounts(stoppinggc, p.histfactor, p.counts[:])
}

func stoppingPausesGCLayout(samples []metrics.Sample) Heatmap {
	idxstoppinggc := metricIdx["/sched/pauses/stopping/gc:seconds"]

	stoppinggc := samples[idxstoppinggc].Value.Float64Histogram()
	histfactor := downsampleFactor(len(stoppinggc.Buckets), maxBuckets)
	buckets := downsampleBuckets(stoppinggc, histfactor)

	return Heatmap{
		Name:       "TODO(set later)",
		Title:      "Stop-the-world Stopping Latencies (GC)",
		Type:       "heatmap",
		UpdateFreq: 5,
		Colorscale: PinkShades,
		Buckets:    floatseq(len(buckets)),
		CustomData: buckets,
		Hover: HeapmapHover{
			YName: "stopping duration",
			YUnit: "duration",
			ZName: "pauses",
		},
		Layout: HeatmapLayout{
			YAxis: HeatmapYaxis{
				Title:    "stopping duration",
				TickMode: "array",
				TickVals: []float64{6, 13, 20, 26, 33, 39.5, 46, 53, 60, 66, 73, 79, 86},
				TickText: []float64{1e-7, 1e-6, 1e-5, 1e-4, 1e-3, 5e-3, 1e-2, 5e-2, 1e-1, 5e-1, 1, 5, 10},
			},
		},
		InfoText: `This heatmap shows the distribution of individual <b>GC-related</b> stop-the-world <i>stopping latencies</i>.
This is the time it takes from deciding to stop the world until all Ps are stopped.
During this time, some threads may be executing.
Uses <b>/sched/pauses/stopping/gc:seconds</b>.`,
	}
}
