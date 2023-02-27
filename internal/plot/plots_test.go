package plot

import (
	"fmt"
	"runtime/metrics"
	"testing"
)

func TestUnusedRuntimeMetrics(t *testing.T) {
	// This test checks if all runtime metrics are plotted. It can fail, it's
	// just informative.
	for _, m := range metrics.All() {
		if _, ok := usedMetrics[m.Name]; !ok {
			fmt.Printf("runtime/metric %q is not used", m.Name)
		}
	}
}
