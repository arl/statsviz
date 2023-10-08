package plot

import (
	"fmt"
	"runtime/metrics"
	"strings"
	"testing"
)

func TestUnusedRuntimeMetrics(t *testing.T) {
	// This "test" can't fail. It just prints which of the metrics exported by
	// /runtime/metrics are not used in any Statsviz plot.
	for _, m := range metrics.All() {
		if _, ok := usedMetrics[m.Name]; !ok && !strings.HasPrefix(m.Name, "/godebug/") {
			fmt.Printf("runtime/metric %q is not used\n", m.Name)
		}
	}
}
