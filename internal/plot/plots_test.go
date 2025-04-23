package plot

import (
	"fmt"
	"os"
	"runtime/metrics"
	"strings"
	"testing"
	"text/tabwriter"
)

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

func TestUnusedRuntimeMetrics(t *testing.T) {
	// This "test" can't fail. It just prints which of the metrics exported by
	// runtime/metrics are not used in any Statsviz plot.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, m := range metrics.All() {
		if _, ok := usedMetrics[m.Name]; !ok && !strings.HasPrefix(m.Name, "/godebug/") {
			fmt.Fprintf(w, "%s\t%s\n", m.Name, kindstr(m.Kind))
		}
	}
	w.Flush()
}
