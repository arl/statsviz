package statsviz

import (
	"errors"
	"fmt"

	"github.com/arl/statsviz/internal/plot"
)

type TimeSeriesType string

const (
	Scatter TimeSeriesType = "scatter" // Scatter is a time series plot made of lines.
	Bar     TimeSeriesType = "bar"     // Bar is a time series plot made of bars.
)

var (
	ErrNoTimeSeries  = errors.New("user plot must have at least one time series")
	ErrEmptyPlotName = errors.New("user plot name can't be empty")
)

type ErrReservedPlotName string

func (e ErrReservedPlotName) Error() string {
	return fmt.Sprintf("%q is a reserved plot name", string(e))
}

// TODO(arl) comment all fields
type TimeSeries struct {
	Name    string
	Unitfmt string
	HoverOn string
	Value   func() float64
}

type TimeSeriesPlotConfig struct {
	Name       string         // Name is the plot name, it's mandatory and must be unique.
	Title      string         // Title is the plot title, shown above the plot.
	Type       TimeSeriesType // Type is either scatter or bar.
	InfoText   string         // Tooltip is the html-aware text shown when the user clicks on the plot Info icon.
	YAxisTitle string         // YAxisTitle is the title of Y axis.

	//  TODO(arl) add link to d3 format page
	YAxisTickSuffix string       // YAxisTickSuffix is the suffix added to tick values.
	Series          []TimeSeries // Series contains the time series shown on this plot, there must be at least one.
}

func (p TimeSeriesPlotConfig) Build() (TimeSeriesPlot, error) {
	var zero TimeSeriesPlot
	if p.Name == "" {
		return zero, ErrEmptyPlotName
	}
	if plot.IsReservedPlotName(p.Name) {
		return zero, ErrReservedPlotName(p.Name)
	}
	if len(p.Series) == 0 {
		return zero, ErrNoTimeSeries
	}

	var (
		subplots []plot.Subplot
		funcs    []func() float64
	)
	for _, ts := range p.Series {
		subplots = append(subplots, plot.Subplot{
			Name:    ts.Name,
			Unitfmt: ts.Unitfmt,
			HoverOn: ts.HoverOn,
		})
		funcs = append(funcs, ts.Value)
	}

	return TimeSeriesPlot{
		timeseries: &plot.ScatterUserPlot{
			Plot: plot.Scatter{
				Name:     p.Name,
				Title:    p.Title,
				Type:     string(p.Type),
				InfoText: p.InfoText,
				Layout: plot.ScatterLayout{
					Yaxis: plot.ScatterYAxis{
						Title:      p.YAxisTitle,
						TickSuffix: p.YAxisTickSuffix,
					},
				},
				Subplots: subplots,
			},
			Funcs: funcs,
		},
	}, nil
}

// TimeSeriesPlot is an opaque type representing a timeseries plot.
type TimeSeriesPlot struct {
	timeseries *plot.ScatterUserPlot
}
