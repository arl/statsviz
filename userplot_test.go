package statsviz

import (
	"errors"
	"testing"
)

func TestTimeSeriesBuilderErrors(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		tsb := NewTimeSeriesPlot("").
			AddSeries(TimeSeries{}, func() float64 { return 0 })

		if _, err := tsb.Build(); !errors.Is(err, ErrEmptyPlotName) {
			t.Errorf("Build() returned err = %v, want %v", err, ErrEmptyPlotName)
		}
	})
	t.Run("reserved name", func(t *testing.T) {
		tsb := NewTimeSeriesPlot("timestamp").
			AddSeries(TimeSeries{}, func() float64 { return 0 })

		var target ErrReservedPlotName
		if _, err := tsb.Build(); !errors.Is(err, target) {
			t.Errorf("Build() returned err = %v, want %v", err, ErrEmptyPlotName)
		}
	})

	/*
	   	t.Run("empty name", func(t *testing.T) {
	   		plot, err := NewTimeSeriesPlot().
	   			Title("User plot").
	   			Type(Scatter).
	   			Tooltip("some <b>tooltip</b>").
	   			YAxisTitle("objects").
	   			YAxisTickSuffix("B").
	   			AddSeries(TimeSeries{
	   				Name:       "name",
	   				Unitfmt:    "%{y:.4s}B",
	   				HoverOn:    "points+fills",
	   				StackGroup: "one",
	   			}, func() float64 { return 0 }).
	   			Build()
	   		if err != nil {
	   			log.Fatalf("failed to build plot: %v", err)
	   		}

	   })
	*/
}
