package plot

import "testing"

func Test_hasDuplicatePlotNames(t *testing.T) {
	tests := []struct {
		name  string
		plots []UserPlot
		want  string
	}{
		{
			"nil",
			nil,
			"",
		},
		{
			"empty",
			[]UserPlot{},
			"",
		},
		{
			"single scatter",
			[]UserPlot{
				{Scatter: &ScatterUserPlot{Plot: Scatter{Name: "a"}}},
			},
			"",
		},
		{
			"single heatmap",
			[]UserPlot{
				{Heatmap: &HeatmapUserPlot{Plot: Heatmap{Name: "a"}}},
			},
			"",
		},
		{
			"two scatter",
			[]UserPlot{
				{Scatter: &ScatterUserPlot{Plot: Scatter{Name: "a"}}},
				{Heatmap: &HeatmapUserPlot{Plot: Heatmap{Name: "b"}}},
				{Scatter: &ScatterUserPlot{Plot: Scatter{Name: "a"}}},
			},
			"a",
		},
		{
			"two heatmap",
			[]UserPlot{
				{Heatmap: &HeatmapUserPlot{Plot: Heatmap{Name: "a"}}},
				{Scatter: &ScatterUserPlot{Plot: Scatter{Name: "b"}}},
				{Heatmap: &HeatmapUserPlot{Plot: Heatmap{Name: "a"}}},
			},
			"a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasDuplicatePlotNames(tt.plots); got != tt.want {
				t.Errorf("hasDuplicatePlotNames() = %q, want %q", got, tt.want)
			}
		})
	}
}
