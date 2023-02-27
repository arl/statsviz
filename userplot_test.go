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
		if _, err := tsb.Build(); !errors.As(err, &target) {
			t.Errorf("Build() returned err = %v, want %v", err, target)
		}
	})
	t.Run("no time series", func(t *testing.T) {
		tsb := NewTimeSeriesPlot("some name")

		if _, err := tsb.Build(); !errors.Is(err, ErrNoTimeSeries) {
			t.Errorf("Build() returned err = %v, want %v", err, ErrEmptyPlotName)
		}
	})
}
