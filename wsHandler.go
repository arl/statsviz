package statsviz

import (
	"net/http"
	"runtime"
	"time"

	"github.com/arl/statsviz/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Ws upgrades the HTTP server connection to the WebSocket protocol and sends
// application statistics every second.
//
// If the upgrade fails, an HTTP error response is sent to the client.
// The package initialization registers it as /debug/statsviz/ws.
func Ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	// Explicitly ignore this error. We don't want to spam standard output
	// each time the other end of the websocket connection closes.
	_ = sendStats(ws)
}

const defaultSendPeriod = time.Second

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn) error {
	tick := time.NewTicker(defaultSendPeriod)
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
