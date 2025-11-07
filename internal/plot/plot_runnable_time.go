package plot

import "runtime/metrics"

var _ = register(description{
	name: "runnable-time",
	tags: []tag{tagScheduler},
	metrics: []string{
		"/sched/latencies:seconds",
	},
	layout: func(samples []metrics.Sample) Heatmap {
		idxschedlat := metricIdx["/sched/latencies:seconds"]

		schedlat := samples[idxschedlat].Value.Float64Histogram()
		histfactor := downsampleFactor(len(schedlat.Buckets), maxBuckets)
		buckets := downsampleBuckets(schedlat, histfactor)

		return Heatmap{
			Name:       "TODO(set later)",
			Title:      "Time Goroutines Spend in 'Runnable' state",
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
			Layout: HeatmapLayout{
				YAxis: HeatmapYaxis{
					Title:    "duration",
					TickMode: "array",
					TickVals: []float64{6, 13, 20, 26, 33, 39.5, 46, 53, 60, 66, 73, 79, 86},
					TickText: []float64{1e-7, 1e-6, 1e-5, 1e-4, 1e-3, 5e-3, 1e-2, 5e-2, 1e-1, 5e-1, 1, 5, 10},
				},
			},
			InfoText: `This heatmap shows the distribution of the time goroutines have spent in the scheduler in a runnable state before actually running, uses <b>/sched/latencies:seconds</b>.`,
		}
	},
	make: func(indices ...int) metricsGetter {
		return &runnableTime{
			idxschedlat: indices[0],
		}
	},
})

type runnableTime struct {
	histfactor int
	counts     [maxBuckets]uint64

	idxschedlat int
}

func (p *runnableTime) values(samples []metrics.Sample) any {
	if p.histfactor == 0 {
		schedlat := samples[p.idxschedlat].Value.Float64Histogram()
		p.histfactor = downsampleFactor(len(schedlat.Buckets), maxBuckets)
	}

	schedlat := samples[p.idxschedlat].Value.Float64Histogram()

	return downsampleCounts(schedlat, p.histfactor, p.counts[:])
}
