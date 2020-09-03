package statsviz

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/arl/statsviz/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type stats struct {
	Mem          runtime.MemStats
	NumGoroutine int
}

const (
	defaultSendPeriod = time.Second
)

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn) error {
	tick := time.NewTicker(defaultSendPeriod)

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

// Ws upgrades the HTTP server connection to the WebSocket protocol and sends
// application statistics every second.
//
// If the upgrade fails, an HTTP error response is sent to the client.
// The package initialization registers it as /debug/statsviz/ws.
func Ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("can't upgrade HTTP connection to Websocket protocol:", err)
		return
	}
	defer ws.Close()

	err = sendStats(ws)
	if err != nil {
		log.Println(err)
	}
}

// Index responds to a request for /debug/statsviz with the statsviz HTML page
// which shows a live visualization of the statistics sent by the application
// over the websocket handler Ws.
//
// The package initialization registers it as /debug/statsviz/.
var Index = http.StripPrefix("/debug/statsviz/", http.FileServer(assets))

func setupRoutes() {
	http.Handle("/debug/statsviz/", Index)
	http.HandleFunc("/debug/statsviz/ws", Ws)
}

func init() {
	setupRoutes()
}
