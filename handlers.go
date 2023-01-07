package statsviz

import (
	"net/http"
	"strings"
	"time"

	"github.com/arl/statsviz/internal/static"

	"github.com/gorilla/websocket"
)

// Index responds to a request for /debug/statsviz with the statsviz HTML page
// which shows a live visualization of the statistics sent by the application
// over the websocket handler Ws.
var Index = IndexAtRoot(defaultRoot)

// IndexAtRoot returns an index statsviz handler rooted at root. It's useful if
// you desire your server to responds with the statsviz HTML page at a
// path that is different than /debug/statsviz.
func IndexAtRoot(root string) http.HandlerFunc {
	prefix := strings.TrimRight(root, "/") + "/"
	assetsFS := http.FileServer(http.FS(static.Assets))
	return http.StripPrefix(prefix, hijack(assetsFS)).ServeHTTP
}

// Ws is a default Websocket handler, created with NewWsHandler, sending statistics
// at the default frequency of 1 message per second.
var Ws = NewWsHandler(defaultSendFrequency)

// NewWsHandler returns a handler that upgrades the HTTP server connection to the WebSocket
// protocol and sends application statistics at the given frequency.
//
// If the upgrade fails, an HTTP error response is sent to the client.
func NewWsHandler(frequency time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer ws.Close()

		// Explicitly ignore this error. We don't want to spam standard output
		// each time the other end of the websocket connection closes.
		_ = sendStats(ws, frequency)
	}
}
