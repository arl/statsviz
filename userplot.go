package statsviz

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/arl/statsviz/internal/plot"
)

type TimeSeriesType string

const (
	Scatter TimeSeriesType = "scatter"
	Bar     TimeSeriesType = "bar"
)

// GetValueFunc reads the value of a metrics and returns it.
type GetValueFunc = plot.GetValueFunc // func() float64

var (
	ErrNoTimeSeries  = errors.New("user plot must have at least one time series")
	ErrEmptyPlotName = errors.New("user plot name must not be empty")
)

type ErrReservedPlotName string

func (e ErrReservedPlotName) Error() string {
	return fmt.Sprintf("%q is a reserved plot name", string(e))
}

type TimeSeriesBuilder struct {
	s     plot.Scatter
	funcs []GetValueFunc // one func per time series
}

func NewTimeSeriesPlot(name string) *TimeSeriesBuilder {
	return &TimeSeriesBuilder{s: plot.Scatter{Name: name}}
}

func (p *TimeSeriesBuilder) Title(title string) *TimeSeriesBuilder {
	p.s.Title = title
	return p
}

func (p *TimeSeriesBuilder) Type(typ TimeSeriesType) *TimeSeriesBuilder {
	p.s.Type = string(typ)
	return p
}

func (p *TimeSeriesBuilder) Tooltip(tooltip string) *TimeSeriesBuilder {
	p.s.InfoText = tooltip
	return p
}

func (p *TimeSeriesBuilder) YAxisTitle(title string) *TimeSeriesBuilder {
	p.s.Layout.Yaxis.Title = title
	return p
}

func (p *TimeSeriesBuilder) YAxisTickSuffix(suffix string) *TimeSeriesBuilder {
	p.s.Layout.Yaxis.TickSuffix = suffix
	return p
}

type TimeSeries struct {
	Name       string
	Unitfmt    string
	StackGroup string
	HoverOn    string
	Color      color.RGBA
}

// AddSeries adds a time series to the current plot. Plots should hold at least
// one time series.
func (p *TimeSeriesBuilder) AddSeries(ts TimeSeries, getval GetValueFunc) *TimeSeriesBuilder {
	p.s.Subplots = append(p.s.Subplots, plot.Subplot{
		Name:       ts.Name,
		Unitfmt:    ts.Unitfmt,
		StackGroup: ts.StackGroup,
		HoverOn:    ts.HoverOn,
		Color:      plot.RGBString(ts.Color.R, ts.Color.G, ts.Color.B),
	})
	p.funcs = append(p.funcs, getval)
	return p
}

func (p *TimeSeriesBuilder) Build() (UserPlot, error) {
	if p.s.Name == "" {
		return UserPlot{}, ErrEmptyPlotName
	}
	if plot.IsReservedPlotName(p.s.Name) {
		return UserPlot{}, ErrReservedPlotName(p.s.Name)
	}
	if len(p.s.Subplots) == 0 {
		return UserPlot{}, ErrNoTimeSeries
	}

	up := UserPlot{
		scatter: &plot.ScatterUserPlot{
			Plot:  plot.Scatter(p.s),
			Funcs: p.funcs,
		},
	}
	return up, nil
}

// Plot is statsviz user plot.
type UserPlot struct {
	// opaque type only made to be passed from Build to WithUserPlot
	scatter *plot.ScatterUserPlot
	heatmap *plot.HeatmapUserPlot
}
