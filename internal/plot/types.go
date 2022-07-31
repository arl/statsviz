package plot

type (
	PlotsDefinition struct {
		Events []string      `json:"events"`
		Series []interface{} `json:"series"`
	}

	ScatterPlotLayout struct {
		Yaxis ScatterPlotLayoutYAxis `json:"yaxis"`
	}

	ScatterPlotLayoutYAxis struct {
		Title      string `json:"title"`
		TickSuffix string `json:"ticksuffix"`
	}

	ScatterPlotSubplot struct {
		Name       string `json:"name"`
		Unitfmt    string `json:"unitfmt"`
		StackGroup string `json:"stackgroup"`
		HoverOn    string `json:"hoveron"`
		Color      Color  `json:"color"`
	}

	ScatterPlot struct {
		Name       string               `json:"name"`
		Title      string               `json:"title"`
		Type       string               `json:"type"`
		UpdateFreq int                  `json:"updateFreq"`
		HorzEvents string               `json:"horzEvents"`
		Layout     ScatterPlotLayout    `json:"layout"`
		Subplots   []ScatterPlotSubplot `json:"subplots"`
	}

	HeatmapPlot struct {
		Name       string            `json:"name"`
		Title      string            `json:"title"`
		Type       string            `json:"type"`
		UpdateFreq int               `json:"updateFreq"`
		HorzEvents string            `json:"horzEvents"`
		Layout     HeatmapPlotLayout `json:"layout"`
		Heatmap    Heatmap           `json:"heatmap"`
	}

	HeatmapPlotLayout struct {
		Yaxis HeatmapPlotLayoutYAxis `json:"yaxis"`
	}

	HeatmapPlotLayoutYAxis struct {
		Title string `json:"title"`
	}

	Heatmap struct {
		Colorscale []WeightedColor `json:"colorscale"`
		Buckets    []float64       `json:"buckets"`
		CustomData []float64       `json:"custom_data"`
		Hover      HeapmapHover    `json:"hover"`
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
