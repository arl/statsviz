package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func scatterPlot() statsviz.TimeSeriesPlot {
	start := time.Now()

	// Describe a second time series.
	sin := statsviz.TimeSeries{
		Name:    "short sin",
		Unitfmt: "%{y:.4s}B",
		Value: func() float64 {
			val := float64(time.Since(start)) / float64(2*time.Second)
			return math.Sin(val)
		},
	}

	// Build a new plot, showing both time series at once.
	plot, err := statsviz.NewTimeSeriesPlot("sin").
		Title("Sinusoide user plot").
		Type(statsviz.Scatter).
		Tooltip(`This tooltip describe the plot that shows a nice sinusoide time series.<br>
This accepts html tags like <b>bold</b> and <i>italic</i>`).
		YAxisTitle("y unit").
		AddTimeSeries(sin).
		Build()

	if err != nil {
		log.Fatalf("failed to build timeseries plot: %v", err)
	}

	return plot
}

func barPlot() statsviz.TimeSeriesPlot {
	// Describe the 'user logins' time series.
	logins := statsviz.TimeSeries{
		Name:    "user log-in",
		Unitfmt: "%{y:.4s}",
		Value: func() float64 {
			return 1_000*rand.Float64() + 2_000
		},
	}

	// Describe the 'user signins' time series.
	signins := statsviz.TimeSeries{
		Name:    "user sign-in",
		Unitfmt: "%{y:.4s}",
		Value: func() float64 {
			return 100*rand.Float64() + 150
		},
	}

	// Build a new plot, showing both time series at once.
	plot, err := statsviz.NewTimeSeriesPlot("users").
		Title("Users").
		Type(statsviz.Bar).
		Tooltip("This plot shows, in real-time, the count of <b>user logged in</b> and <b>user registered</b> events.").
		YAxisTitle("users").
		AddTimeSeries(logins).
		AddTimeSeries(signins).
		Build()

	if err != nil {
		log.Fatalf("failed to build timeseries plot: %v", err)
	}

	return plot
}

func main() {
	go example.Work()

	mux := http.NewServeMux()

	// Create statsviz endpoint with 2 additional, user-provided plots.
	ep := statsviz.NewEndpoint(
		statsviz.WithTimeseriesPlot(scatterPlot()),
		statsviz.WithTimeseriesPlot(barPlot()),
	)

	// Register the endpoint handlers on the mux.
	ep.Register(mux)

	fmt.Println("Point your browser to http://localhost:8093/debug/statsviz/")
	log.Fatal(http.ListenAndServe(":8093", mux))
}
