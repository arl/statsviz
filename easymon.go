package rtprof

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
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

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
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

func setupRoutes() {
	http.Handle("/debug/rtprof/", http.StripPrefix("/debug/rtprof/", http.FileServer(assets)))
	http.HandleFunc("/debug/rtprof/ws", wsEndpoint)
}

func init() {
	setupRoutes()
}
