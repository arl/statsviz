package plot

type (
	Definition struct {
		Events []string      `json:"events"`
		Series []interface{} `json:"series"`
	}

	ScatterLayout struct {
		Yaxis ScatterLayoutYAxis `json:"yaxis"`
	}

	ScatterLayoutYAxis struct {
		Title      string `json:"title"`
		TickSuffix string `json:"ticksuffix"`
	}

	Subplot struct {
		Name       string `json:"name"`
		Unitfmt    string `json:"unitfmt"`
		StackGroup string `json:"stackgroup"`
		HoverOn    string `json:"hoveron"`
		Color      Color  `json:"color"`
	}

	Scatter struct {
		Name       string        `json:"name"`
		Title      string        `json:"title"`
		Type       string        `json:"type"`
		UpdateFreq int           `json:"updateFreq"`
		HorzEvents string        `json:"horzEvents"`
		Layout     ScatterLayout `json:"layout"`
		Subplots   []Subplot     `json:"subplots"`
	}

	HeatmapPlot struct {
		Name       string            `json:"name"`
		Title      string            `json:"title"`
		Type       string            `json:"type"`
		UpdateFreq int               `json:"updateFreq"`
		HorzEvents string            `json:"horzEvents"`
		Layout     HeatmapPlotLayout `json:"layout"`
		Colorscale []WeightedColor   `json:"colorscale"`
		Buckets    []float64         `json:"buckets"`
		CustomData []float64         `json:"custom_data"`
		Hover      HeapmapHover      `json:"hover"`
	}

	HeatmapPlotLayout struct {
		Yaxis HeatmapPlotLayoutYAxis `json:"yaxis"`
	}

	HeatmapPlotLayoutYAxis struct {
		Title string `json:"title"`
	}

	HeapmapHover struct {
		YName string `json:"yname"`
		YUnit string `json:"yunit"` // 'duration', 'bytes' or custom
		ZName string `json:"zname"`
	}

	Axis struct {
		Title string
		Unit  Unit
	}

	Unit struct {
		TickSuffix string
		UnitFmt    string
	}
)

var Bytes = Unit{TickSuffix: "B", UnitFmt: "%{y:.4s}B"}
