package statsviz

import (
	"time"

	"github.com/gorilla/websocket"
)

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	// If the websocket connection is initiated by an already open web ui
	// (started by a previous process for example) then plotsdef.js won't be
	// requested. So, call plots.config manually to ensure that the data
	// structures inside 'plots' are correctly initialized.
	plots.config()

	for range tick.C {
		if err := conn.WriteJSON(plots.values()); err != nil {
			return err
		}
	}

	panic("unreachable")
}
