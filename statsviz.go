package statsviz

import (
	"encoding/json"
	"expvar"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// sendStats indefinitely send runtime statistics on the websocket connection.
func sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	expvar.Publish("NumGoroutine", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	for {
		select {
		case <-tick.C:
			var s string
			expvar.Do(func(kv expvar.KeyValue) {
				s += fmt.Sprintf("\"%v\": %v,", kv.Key, kv.Value.String())
			})
			sj, _ := json.Marshal(s)
			j := strings.ReplaceAll(string(sj[1:len(sj)-1]), "\\", "")
			stats := fmt.Sprintf("{%v\"null\":false}", j)
			if err := conn.WriteJSON(stats); err != nil {
				return err
			}
		}

	}
}
