package statsviz

import (
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

type stats struct {
	GoVersion    string
	Mem          runtime.MemStats
	NumGoroutine int
	CustomData   interface{}
}
type CustomDataGenerateFunc func() interface{}

var CustomDataGenerate CustomDataGenerateFunc

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	stats := stats{GoVersion: runtime.Version()}
	for range tick.C {
		runtime.ReadMemStats(&stats.Mem)
		stats.NumGoroutine = runtime.NumGoroutine()
		if CustomDataGenerate != nil {
			stats.CustomData = CustomDataGenerate()
		} else {
			stats.CustomData = nil
		}
		if err := conn.WriteJSON(stats); err != nil {
			return err
		}
	}

	panic("unreachable")
}
