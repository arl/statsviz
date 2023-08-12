package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func scatterPlot() statsviz.TimeSeriesPlot {
	val := 0.

	// Describe the 'sine' time series.
	sine := statsviz.TimeSeries{
		Name:    "short sin",
		Unitfmt: "%{y:.4s}B",
		GetValue: func() float64 {
			val += 0.5
			return math.Sin(val)
		},
	}

	// Build a new plot, showing our sine time series
	plot, err := statsviz.TimeSeriesPlotConfig{
		Name:  "sine",
		Title: "Sine",
		Type:  statsviz.Scatter,
		InfoText: `This tooltip describe the plot that shows a <i>sine</i> time series.<br>
This accepts HTML tags like <b>bold</b> and <i>italic</i>`,
		YAxisTitle: "y unit",
		Series:     []statsviz.TimeSeries{sine},
	}.Build()
	if err != nil {
		log.Fatalf("failed to build timeseries plot: %v", err)
	}

	return plot
}

func barPlot() statsviz.TimeSeriesPlot {
	// Describe the 'user logins' time series.
	logins := statsviz.TimeSeries{
		Name:    "user logins",
		Unitfmt: "%{y:.4s}",
		GetValue: func() float64 {
			return 1_000*rand.Float64() + 2_000
		},
	}

	// Describe the 'user signins' time series.
	signins := statsviz.TimeSeries{
		Name:    "user signins",
		Unitfmt: "%{y:.4s}",
		GetValue: func() float64 {
			return 100*rand.Float64() + 150
		},
	}

	// Build a new plot, showing both time series at once.
	plot, err := statsviz.TimeSeriesPlotConfig{
		Name:  "users",
		Title: "Users",
		Type:  statsviz.Bar,
		InfoText: `This plot shows the real time count of <i>user login</i> and <i>user signin</i> events.<br>
This accepts HTML tags like <b>bold</b> and <i>italic</i>`,
		YAxisTitle: "users",
		Series:     []statsviz.TimeSeries{logins, signins},
	}.Build()
	if err != nil {
		log.Fatalf("failed to build timeseries plot: %v", err)
	}

	return plot
}

func main() {
	go example.Work()

	mux := http.NewServeMux()

	// Register statsviz handlers with 2 additional plots, user-provided plots.
	_ = statsviz.Register(mux,
		statsviz.TimeseriesPlot(scatterPlot()),
		statsviz.TimeseriesPlot(barPlot()),
	)

	fmt.Println("Point your browser to http://localhost:8093/debug/statsviz/")
	log.Fatal(http.ListenAndServe(":8093", mux))
}
