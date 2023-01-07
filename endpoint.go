package statsviz

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arl/statsviz/internal/static"
	"github.com/gorilla/websocket"
)

type Endpoint struct {
	intv time.Duration
	root string
}

func NewEndpoint() *Endpoint {
	const (
		defaultRoot         = "/debug/statsviz"
		defaultSendInterval = time.Second
	)

	return &Endpoint{
		intv: defaultSendInterval,
		root: defaultRoot,
	}
}

// WithSendInterval specifies the time interval at which metrics are requested
// and sent to the user interface. Default is one second.
func (e *Endpoint) WithSendInterval(intv time.Duration) *Endpoint {
	e.intv = intv
	return e
}

// WithRoot specifies the root path of statsviz HTTP handlers.
// Default is /debug/statsviz.
func (e *Endpoint) WithRoot(root string) *Endpoint {
	e.root = root
	return e
}

// indexAtRoot returns an index statsviz handler rooted at root. It's useful if
// you desire your server to responds with the statsviz HTML page at a
// path that is different than /debug/statsviz.
func indexAtRoot(root string) http.HandlerFunc {
	prefix := strings.TrimRight(root, "/") + "/"
	assetsFS := http.FileServer(http.FS(static.Assets))
	return http.StripPrefix(prefix, hijack(assetsFS)).ServeHTTP
}

// Register registers on the given mux the HTTP handlers required for statsviz
// endpoint.
func (e *Endpoint) Register(mux *http.ServeMux) {
	mux.Handle(e.root+"/", e.Index())
	mux.HandleFunc(e.root+"/ws", e.Ws())
}

// Index returns the index handler, responding with statsviz user interface HTML
// page. Use [WithRoot] if you wish statsviz user interface to be presented at
// path different than /debug/statsviz.
func (e *Endpoint) Index() http.HandlerFunc {
	prefix := strings.TrimSuffix(e.root, "/") + "/"
	assetsFS := http.FileServer(http.FS(static.Assets))
	return http.StripPrefix(prefix, hijack(assetsFS)).ServeHTTP
}

// Ws returns a handler that upgrades the HTTP connection to the WebSocket
// protocol and sends application statistics.
func (e *Endpoint) Ws() http.HandlerFunc {
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
		_ = sendStats(ws, e.intv)
	}
}

// hijack returns a handler that hijacks requests for plotsdef.js, this file is
// generated dynamically. Other requests are forwarded to h, typically a http
// file server.
func hijack(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "js/plotsdef.js" {
			buf := &bytes.Buffer{}
			buf.WriteString("export default ")
			enc := json.NewEncoder(buf)
			enc.SetIndent("", "  ")
			if err := enc.Encode(plots.Config()); err != nil {
				panic("error encoding plots definition: " + err.Error())
			}
			buf.WriteString(";")
			w.Header().Add("Content-Length", strconv.Itoa(buf.Len()))
			w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
			buf.WriteTo(w)
			return
		}
		// Force Content-Type if needed.
		if ct, ok := contentTypes[r.URL.Path]; ok {
			w.Header().Add("Content-Type", ct)
		}

		h.ServeHTTP(w, r)
	}
}

// Force Content-Type HTTP header for certain files of some javascript libraries
// that have no extensions. Otherwise the http fileserver would serve them under
// "Content-Type = text/plain".
var contentTypes = map[string]string{
	"libs/js/popperjs-core2": "text/javascript",
	"libs/js/tippy.js@6":     "text/javascript",
}
