// Package statsviz allows to visualise Go program runtime metrics data in real
// time: heap, objects, goroutines, GC pauses, scheduler, etc. in your browser.
//
// Create a statsviz [Endpoint] and register it with your server [http.ServeMux]
// (preferred method):
//
//	mux := http.NewServeMux()
//	endpoint := statviz.NewEndpoint()
//	endpoint.Register(mux)
//
// Or register with [http.DefaultServeMux`]:
//
//	endpoint := statviz.NewEndpoint()
//	endpoint.Register(http.DefaultServeMux)
//
// By default Statsviz is served at `/debug/statsviz/`. You can change that (and
// other things) using methods on the [statsviz.Endpoint] instance.
//
// If your application is not already running an HTTP server, you need to start
// one. Add "net/http" and "log" to your imports and the following code to your
// main function:
//
//	go func() {
//	    log.Println(http.ListenAndServe("localhost:6060", nil))
//	}()
//
// Then open your browser at http://localhost:6060/debug/statsviz/.
package statsviz

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arl/statsviz/internal/plot"
	"github.com/arl/statsviz/internal/static"

	"github.com/gorilla/websocket"
)

const (
	defaultRoot         = "/debug/statsviz"
	defaultSendInterval = time.Second
)

// An Endpoint serves, and consists of, 2 HTTP handlers necessary for Statsviz
// user interface.
type Endpoint struct {
	intv  time.Duration // interval between consecutive metrics emission
	root  string        // http path root
	plots *plot.List    // plots shown on the user interface

	// user plots
	userPlots []plot.UserPlot
}

// NewEndpoint constructs a new Statsviz endpoint, pre-configured with default
// settings or with given options.
func NewEndpoint(opts ...Option) *Endpoint {
	e := &Endpoint{
		intv: defaultSendInterval,
		root: defaultRoot,
	}

	for _, opt := range opts {
		opt(e)
	}

	e.plots = plot.NewList(e.userPlots)
	return e
}

// Option is an Endpoint configuration option.
type Option func(*Endpoint)

// WithInterval option changes the interval between successive acquisitions of
// metrics and their sending to the user interface. By default, the interval is
// one second.
func WithInterval(intv time.Duration) Option {
	return func(e *Endpoint) {
		e.intv = intv
	}
}

// WithRoot option changes the root path at which statsviz endpoint is served on
// the HTTP server. By default this path is /debug/statviz.
func WithRoot(path string) Option {
	return func(e *Endpoint) {
		e.root = path
	}
}

// WithTimeseriesPlot adds a new timeseries plot to statsviz ui that tracks some
// user-provided metric. Can be called multiple times.
func WithTimeseriesPlot(tsp TimeSeriesPlot) Option {
	return func(e *Endpoint) {
		e.userPlots = append(e.userPlots,
			plot.UserPlot{Scatter: tsp.timeseries})
	}
}

// Register registers statviz HTTP handlers on the provided mux.
func (e *Endpoint) Register(mux *http.ServeMux) {
	mux.Handle(e.root+"/", e.Index())
	mux.HandleFunc(e.root+"/ws", e.Ws())
}

// Index returns the index handler, responding with statsviz user interface HTML
// page. By default, the returned handler is served at /debug/statsviz. Use
// [WithRoot] to change that path.
func (e *Endpoint) Index() http.HandlerFunc {
	prefix := strings.TrimSuffix(e.root, "/") + "/"
	assetsFS := http.FileServer(http.FS(static.Assets))
	return http.StripPrefix(prefix, hijack(assetsFS, e.plots)).ServeHTTP
}

// Ws returns the websocket handler statsviz uses to send application metrics.
//
// The underlying net.Conn is used to upgrade the HTTP server connection to the
// websocket protocol.
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

		// Explicitly ignore this error. This happens if/when the other end
		// connection closes for example. We can't handle it in any meaningful
		// way anyways.
		_ = e.sendStats(ws, e.intv)
	}
}

// hijack returns a handler that hijacks requests for plotsdef.js, this file is
// generated dynamically. Other requests are forwarded to h, typically a http
// file server.
func hijack(h http.Handler, plots *plot.List) http.HandlerFunc {
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

// sendStats indefinitely send runtime statistics on the websocket connection.
func (e *Endpoint) sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	// If the websocket connection is initiated by an already open web ui
	// (started by a previous process for example) then plotsdef.js won't be
	// requested. So, call plots.Config() manually to ensure that the data
	// structures inside 'plots' are correctly initialized.
	e.plots.Config()

	for range tick.C {
		w, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return err
		}
		if err := e.plots.WriteValues(w); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	}

	panic("unreachable")
}
