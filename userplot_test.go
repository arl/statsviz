package statsviz

import (
	"errors"
	"testing"
)

func TestTimeSeriesPlotConfigErrors(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		tsb := TimeSeriesPlotConfig{}
		if _, err := tsb.Build(); !errors.Is(err, ErrEmptyPlotName) {
			t.Errorf("Build() returned err = %v, want %v", err, ErrEmptyPlotName)
		}
	})
	t.Run("reserved name", func(t *testing.T) {
		tsb := TimeSeriesPlotConfig{Name: "timestamp"}
		var target ErrReservedPlotName
		if _, err := tsb.Build(); !errors.As(err, &target) {
			t.Errorf("Build() returned err = %v, want %v", err, target)
		}
	})
	t.Run("missing value function", func(t *testing.T) {
		config := TimeSeriesPlotConfig{
			Name:   "some name",
			Series: []TimeSeries{{Name: "some series"}},
		}
		if _, err := config.Build(); err == nil {
			t.Error("Build() accepted a time series without a GetValue function")
		}
	})
	t.Run("no time series", func(t *testing.T) {
		tsb := TimeSeriesPlotConfig{Name: "some name"}
		if _, err := tsb.Build(); !errors.Is(err, ErrNoTimeSeries) {
			t.Errorf("Build() returned err = %v, want %v", err, ErrEmptyPlotName)
		}
	})
}
