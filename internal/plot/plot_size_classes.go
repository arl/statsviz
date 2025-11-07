package plot

import "runtime/metrics"

type sizeClasses struct {
	sizeClasses []uint64

	idxallocs int
	idxfrees  int
}

func makeSizeClasses(indices ...int) metricsGetter {
	return &sizeClasses{
		idxallocs: indices[0],
		idxfrees:  indices[1],
	}
}

func (p *sizeClasses) values(samples []metrics.Sample) any {
	allocsBySize := samples[p.idxallocs].Value.Float64Histogram()
	freesBySize := samples[p.idxfrees].Value.Float64Histogram()

	if p.sizeClasses == nil {
		p.sizeClasses = make([]uint64, len(allocsBySize.Counts))
	}

	for i := range p.sizeClasses {
		p.sizeClasses[i] = allocsBySize.Counts[i] - freesBySize.Counts[i]
	}
	return p.sizeClasses
}

func sizeClassesLayout(samples []metrics.Sample) Heatmap {
	idxallocs := metricIdx["/gc/heap/allocs-by-size:bytes"]
	idxfrees := metricIdx["/gc/heap/frees-by-size:bytes"]

	// Perform a sanity check on the number of buckets on the 'allocs' and
	// 'frees' size classes histograms. Statsviz plots a single histogram based
	// on those 2 so we want them to have the same number of buckets, which
	// should be true.
	allocsBySize := samples[idxallocs].Value.Float64Histogram()
	freesBySize := samples[idxfrees].Value.Float64Histogram()
	if len(allocsBySize.Buckets) != len(freesBySize.Buckets) {
		panic("different number of buckets in allocs and frees size classes histograms")
	}

	// No downsampling for the size classes histogram (factor=1) but we still
	// need to adapt boundaries for plotly heatmaps.
	buckets := downsampleBuckets(allocsBySize, 1)

	return Heatmap{
		Name:       "TODO(set later)",
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
		Layout: HeatmapLayout{
			YAxis: HeatmapYaxis{
				Title:    "size class",
				TickMode: "array",
				TickVals: []float64{1, 9, 17, 25, 31, 37, 43, 50, 58, 66},
				TickText: []float64{1 << 4, 1 << 7, 1 << 8, 1 << 9, 1 << 10, 1 << 11, 1 << 12, 1 << 13, 1 << 14, 1 << 15},
			},
		},
	}
}
