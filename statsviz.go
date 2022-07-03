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

	for range tick.C {
		metrics.Read(samples)
		if err := conn.WriteJSON(plotsValues(samples)); err != nil {
			return err
		}
	}

	panic("unreachable")
}
