package plot

import "runtime/metrics"

var _ = register(description{
	name: "stopping-pauses-other",
	tags: []tag{tagScheduler},
	metrics: []string{
		"/sched/pauses/stopping/other:seconds",
	},
	layout: func(samples []metrics.Sample) Heatmap {
		idxstoppingother := metricIdx["/sched/pauses/stopping/other:seconds"]

		stoppingother := samples[idxstoppingother].Value.Float64Histogram()
		histfactor := downsampleFactor(len(stoppingother.Buckets), maxBuckets)
		buckets := downsampleBuckets(stoppingother, histfactor)

		return Heatmap{
			Name:       "TODO(set later)",
			Title:      "Stop-the-world Stopping Latencies (Other)",
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
			InfoText: `This heatmap shows the distribution of individual <b>non-GC-related</b> stop-the-world <i>stopping latencies</i>.
This is the time it takes from deciding to stop the world until all Ps are stopped.
This is a subset of the total non-GC-related stop-the-world time. During this time, some threads may be executing.
Uses <b>/sched/pauses/stopping/other:seconds</b>.`,
		}
	},
	make: func(indices ...int) metricsGetter {
		return &stoppingPausesOther{
			idxstoppingother: indices[0],
		}
	},
})

type stoppingPausesOther struct {
	histfactor int
	counts     [maxBuckets]uint64

	idxstoppingother int
}

func (p *stoppingPausesOther) values(samples []metrics.Sample) any {
	if p.histfactor == 0 {
		stoppingother := samples[p.idxstoppingother].Value.Float64Histogram()
		p.histfactor = downsampleFactor(len(stoppingother.Buckets), maxBuckets)
	}

	stoppingother := samples[p.idxstoppingother].Value.Float64Histogram()
	return downsampleCounts(stoppingother, p.histfactor, p.counts[:])
}
