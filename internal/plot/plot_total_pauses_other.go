package plot

import "runtime/metrics"

type totalPausesOther struct {
	histfactor int
	counts     [maxBuckets]uint64

	idxtotalother int
}

func makeTotalPausesOther(indices ...int) metricsGetter {
	return &totalPausesOther{
		idxtotalother: indices[0],
	}
}

func (p *totalPausesOther) values(samples []metrics.Sample) any {
	if p.histfactor == 0 {
		totalother := samples[p.idxtotalother].Value.Float64Histogram()
		p.histfactor = downsampleFactor(len(totalother.Buckets), maxBuckets)
	}

	totalother := samples[p.idxtotalother].Value.Float64Histogram()
	return downsampleCounts(totalother, p.histfactor, p.counts[:])
}

func totalPausesOtherLayout(samples []metrics.Sample) Heatmap {
	idxtotalother := metricIdx["/sched/pauses/total/other:seconds"]

	totalother := samples[idxtotalother].Value.Float64Histogram()
	histfactor := downsampleFactor(len(totalother.Buckets), maxBuckets)
	buckets := downsampleBuckets(totalother, histfactor)

	return Heatmap{
		Name:       "TODO(set later)",
		Title:      "Stop-the-world Pause Latencies (Other)",
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
		Layout: HeatmapLayout{
			YAxis: HeatmapYaxis{
				Title:    "pause duration",
				TickMode: "array",
				TickVals: []float64{6, 13, 20, 26, 33, 39.5, 46, 53, 60, 66, 73, 79, 86},
				TickText: []float64{1e-7, 1e-6, 1e-5, 1e-4, 1e-3, 5e-3, 1e-2, 5e-2, 1e-1, 5e-1, 1, 5, 10},
			},
		},
		InfoText: `This heatmap shows the distribution of individual <b>non-GC-related</b> stop-the-world <i>pause latencies</i>.
This is the time from deciding to stop the world until the world is started again.
Some of this time is spent getting all threads to stop (measured directly in <i>/sched/pauses/stopping/other:seconds</i>).
Uses <b>/sched/pauses/total/other:seconds</b>.`,
	}
}
