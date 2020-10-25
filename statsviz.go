package statsviz

import (
	"runtime"
	"time"

	"github.com/arl/statsviz/websocket"
)

type stats struct {
	Mem          runtime.MemStats
	NumGoroutine int
}

const defaultSendFrequency = time.Second

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn) error {
	tick := time.NewTicker(defaultSendFrequency)
	defer tick.Stop()

	var stats stats
	for {
		select {
		case <-tick.C:
			runtime.ReadMemStats(&stats.Mem)
			stats.NumGoroutine = runtime.NumGoroutine()
			if err := conn.WriteJSON(stats); err != nil {
				return err
			}
		}
	}
}
