package plot

type (
	Config struct {
		// Series contains the plots we want to show and how we want to show them.
		Series []interface{} `json:"series"`
		// Events contains a list of 'events time series' names. Series with
		// these names must be sent alongside other series. An event time series
		// is just made of timestamps with no associated value, each of which
		// gets plotted as a vertical line over another plot.
		Events []string `json:"events"`
	}

	Scatter struct {
		Name       string `json:"name"`
		Title      string `json:"title"`
		Type       string `json:"type"`
		UpdateFreq int    `json:"updateFreq"`
		InfoText   string `json:"infoText"`
		Events     string `json:"events"`
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
		InfoText   string `json:"infoText"`
		Events     string `json:"events"`
		Layout     struct {
			Yaxis struct {
				Title    string    `json:"title"`
				TickMode string    `json:"tickmode"`
				TickVals []float64 `json:"tickvals"`
				TickText []float64 `json:"ticktext"`
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
