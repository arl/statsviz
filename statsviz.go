package statsviz

import (
	"time"

	"github.com/gorilla/websocket"
)

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	for range tick.C {
		if err := conn.WriteJSON(plotsValues()); err != nil {
			return err
		}
	}

	panic("unreachable")
}
