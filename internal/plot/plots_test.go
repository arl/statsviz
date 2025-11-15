package plot

import (
	"fmt"
	"maps"
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
	unused := maps.Clone(reg().allnames)

	reg := reg()
	for _, d := range reg.descriptions {
		for _, m := range d.metrics {
			delete(unused, m)
		}
	}

	// remove godebug metrics
	for m := range unused {
		if strings.HasPrefix(m, "/godebug/") {
			delete(unused, m)
		}
	}

	if len(unused) == 0 {
		return
	}
	t.Log("some runtime metrics are not used by any plot:\n")
	all := metrics.All()
	w := tabwriter.NewWriter(os.Stderr, 0, 8, 2, ' ', 0)
	for m := range unused {
		desc := all[slices.IndexFunc(all, func(desc metrics.Description) bool { return desc.Name == m })]
		fmt.Fprintf(w, "\t%s\t%s\t%s\n", desc.Name, kindstr(desc.Kind), clampstr(desc.Description))
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
