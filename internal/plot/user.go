package plot

type ScatterUserPlot struct {
	Plot  Scatter
	Funcs []func() float64
}

type HeatmapUserPlot struct {
	Plot Heatmap
	// TODO(arl): heatmap get value func
}

type UserPlot struct {
	Scatter *ScatterUserPlot
	Heatmap *HeatmapUserPlot
}

func (up UserPlot) Layout() interface{} {
	switch {
	case (up.Scatter != nil) == (up.Heatmap != nil):
		panic("userplot must be a timeseries or a heatmap")
	case up.Scatter != nil:
		return up.Scatter.Plot
	case up.Heatmap != nil:
		return up.Heatmap.Plot
	}

	panic("unreeachable")
}
