package statsviz

import (
	"runtime/metrics"
	"time"

	"github.com/gorilla/websocket"
)

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	// The web ui may already be running, started from a previous process, in
	// which case it won't try to fetch the plot definitions. Nevertheless on
	// the Go side we need to initialize the plot defs.
	plotsdef()

	for range tick.C {
		metrics.Read(samples)
		if err := conn.WriteJSON(plotsValues(samples)); err != nil {
			return err
		}
	}

	panic("unreachable")
}
