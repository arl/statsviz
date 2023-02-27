package plot

// type GetValueFunc func() float64

type ScatterUserPlot struct {
	Plot  Scatter
	Funcs []func() float64
}

type HeatmapUserPlot struct {
	Plot Heatmap
	// TODO(arl): heatmap get value func
}

type UserPlot struct {
	// TODO(arl) we should create a NewUserPlot constructor so that we can
	// unexport both fields. Once they're unexported, there's no reason to keep
	// 2 structs statsviz.UserPlot and internal.plot.UserPlot. We can simply
	// keep the internal one and do:
	//
	//  package statsviz
	//  type UserPlot plot.UserPlot (this will become opaque and that's what we want!)
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
