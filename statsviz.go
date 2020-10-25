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
