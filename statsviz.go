package statsviz

import (
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

type stats struct {
	Mem          runtime.MemStats
	NumGoroutine int
}

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	var (
		stats stats
		err   error
	)
	for range tick.C {
		runtime.ReadMemStats(&stats.Mem)
		stats.NumGoroutine = runtime.NumGoroutine()
		if err = conn.WriteJSON(stats); err != nil {
			break
		}
	}

	return err
}
