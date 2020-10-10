// Package statsviz serves via its HTTP server an HTML page displaying live
// visualization of the application runtime statistics.
//
// The package is typically only imported for the side effect of
// registering its HTTP handler.
// The handled path is /debug/statsviz/.
//
// To use statsviz, link this package into your program:
//	import _ "github.com/arl/statsviz"
//
// If your application is not already running an http server, you
// need to start one. Add "net/http" and "log" to your imports and
// the following code to your main function:
//
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()
//
// If you are not using DefaultServeMux, you will have to register handlers
// with the mux you are using.
//
// Then open your browser at http://localhost:6060/debug/statsviz/
package statsviz

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/arl/statsviz/websocket"
)

// Register registers statsviz HTTP handlers on the provided mux.
func Register(mux *http.ServeMux) {
	mux.Handle("/debug/statsviz/", Index)
	mux.HandleFunc("/debug/statsviz/ws", Ws)
}

// RegisterDefault registers statsviz HTTP handlers on the default serve mux.
//
// Note this is not advised on a production server, unless it only serves on
// localhost.
func RegisterDefault() {
	Register(http.DefaultServeMux)
}

// Index responds to a request for /debug/statsviz with the statsviz HTML page
// which shows a live visualization of the statistics sent by the application
// over the websocket handler Ws.
//
// The package initialization registers it as /debug/statsviz/.
var Index = http.StripPrefix("/debug/statsviz/", http.FileServer(assets))

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type stats struct {
	Mem          runtime.MemStats
	NumGoroutine int
}

const defaultSendPeriod = time.Second

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
