package statsviz

import (
	"errors"
	"testing"
)

func TestTimeSeriesBuilderErrors(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		tsb := NewTimeSeriesPlot("").
			AddTimeSeries(TimeSeries{})

		if _, err := tsb.Build(); !errors.Is(err, ErrEmptyPlotName) {
			t.Errorf("Build() returned err = %v, want %v", err, ErrEmptyPlotName)
		}
	})
	t.Run("reserved name", func(t *testing.T) {
		tsb := NewTimeSeriesPlot("timestamp").
			AddTimeSeries(TimeSeries{})

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
