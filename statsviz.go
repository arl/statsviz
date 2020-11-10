package statsviz

import (
	"runtime"
	"time"
)

type stats struct {
	Mem          runtime.MemStats
	NumGoroutine int
}

// various websocket connection interface
type conn interface {
	WriteJSON(interface{}) error
}

const defaultSendFrequency = time.Second

// SendStats indefinitely send runtime statistics on the websocket connection.
func SendStats(c conn) error {
	tick := time.NewTicker(defaultSendFrequency)
	defer tick.Stop()

	var stats stats
	for {
		select {
		case <-tick.C:
			runtime.ReadMemStats(&stats.Mem)
			stats.NumGoroutine = runtime.NumGoroutine()
			if err := c.WriteJSON(stats); err != nil {
				return err
			}
		}
	}
}
