package plot

type (
	Config struct {
		// Events are transversal time series, which can be plotted as
		// horizontal lines on any plots.
		Events []string `json:"events"`
		// Series contains the plots we want to show and how we want to show them.
		Series []interface{} `json:"series"`
	}

	Scatter struct {
		Name       string `json:"name"`
		Title      string `json:"title"`
		Type       string `json:"type"`
		UpdateFreq int    `json:"updateFreq"`
		HorzEvents string `json:"horzEvents"`
		Layout     struct {
			Yaxis struct {
				Title      string `json:"title"`
				TickSuffix string `json:"ticksuffix"`
			} `json:"yaxis"`
		} `json:"layout"`
		Subplots []Subplot `json:"subplots"`
	}

	Subplot struct {
		Name       string `json:"name"`
		Unitfmt    string `json:"unitfmt"`
		StackGroup string `json:"stackgroup"`
		HoverOn    string `json:"hoveron"`
		Color      string `json:"color"`
	}

	Heatmap struct {
		Name       string `json:"name"`
		Title      string `json:"title"`
		Type       string `json:"type"`
		UpdateFreq int    `json:"updateFreq"`
		HorzEvents string `json:"horzEvents"`
		Layout     struct {
			Yaxis struct {
				Title string `json:"title"`
			} `json:"yaxis"`
		} `json:"layout"`
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
