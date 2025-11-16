package plot

import (
	"fmt"
	"os"
	"runtime/metrics"
	"slices"
	"strings"
	"testing"
	"text/tabwriter"
)

func TestUnusedRuntimeMetrics(t *testing.T) {
	// This test just prints the metrics we're not using in any plot. It can't
	// fail, it's informational.
	used := make(map[string]bool)
	for _, d := range reg().descriptions {
		for _, m := range d.metrics {
			used[m] = true
		}
	}

	// Discard godebug metrics and used metrics.
	all := metrics.All()
	all = slices.DeleteFunc(all, func(desc metrics.Description) bool {
		return strings.HasPrefix(desc.Name, "/godebug/")
	})
	all = slices.DeleteFunc(all, func(desc metrics.Description) bool {
		return used[desc.Name]
	})

	if len(all) == 0 {
		t.Log("all metrics are used!")
		return
	}

	t.Log("some runtime metrics are not used by any plot:\n")

	w := tabwriter.NewWriter(os.Stderr, 0, 8, 2, ' ', 0)
	for _, m := range all {
		fmt.Fprintf(w, "\t%s\t%s\t%s\n", m.Name, kindstr(m.Kind), clampstr(m.Description))
	}
	w.Flush()
}

func kindstr(k metrics.ValueKind) string {
	switch k {
	case metrics.KindUint64:
		return "uint64"
	case metrics.KindFloat64:
		return "float64"
	case metrics.KindFloat64Histogram:
		return "float64 histogram"
	default:
		return "unknown"
	}
}

func clampstr(s string) string {
	const maxlen = 80
	if len(s) > maxlen {
		return s[:maxlen-3] + "..."
	}
	return s
}
