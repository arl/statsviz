// Package statsviz allows visualizing Go program runtime metrics data in real
// time, including heap, objects, goroutines, GC pauses, scheduler and much
// more, in your browser.
//
// To use Statsviz, create a new Statsviz [Server] and register it with your
// HTTP server's [http.ServeMux] (preferred method):
//
//	mux := http.NewServeMux()
//	ss := statviz.Server{}
//	ss.Register(mux)
//
// Alternatively, you can register with [http.DefaultServeMux]:
//
//	ss := statviz.Server{}
//	s.Register(http.DefaultServeMux)
//
// By default, Statsviz is served at `/debug/statsviz/`. You can change this and
// other settings by passing some [Option] to [NewServer].
//
// If your application is not already running an HTTP server, you need to start
// one. Add "net/http" and "log" to your imports, and use the following code in
// your main function:
//
//	go func() {
//	    log.Println(http.ListenAndServe("localhost:6060", nil))
//	}()
//
// Then open your browser and visit http://localhost:6060/debug/statsviz/.
package statsviz

// TODO(arl) Make sure to keep the README.md file updated with this documentation.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/arl/statsviz/internal/plot"
	"github.com/arl/statsviz/internal/static"
)

const (
	defaultRoot         = "/debug/statsviz"
	defaultSendInterval = time.Second
)

// Server is the core component of Statsviz. It collects and periodically
// updates metrics data and provides two essential HTTP handlers:
//   - the Index handler serves Statsviz user interface, allowing you to
//     visualize runtime metrics on your browser.
//   - The Ws handler establishes a WebSocket connection allowing the connected
//     browser to receive metrics updates from the server.
//
// The zero value if a valid Server.
type Server struct {
	intv      time.Duration // interval between consecutive metrics emission
	root      string        // HTTP path root
	plots     *plot.List    // plots shown on the user interface
	userPlots []plot.UserPlot
}

// NewServer constructs a new Statsviz Server with default settings or given
// options, if any.
func NewServer(opts ...Option) *Server {
	s := &Server{
		intv: defaultSendInterval,
		root: defaultRoot,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.plots = plot.NewList(s.userPlots)
	return s
}

// Option is a configuration option for the Server.
type Option func(*Server)

// WithInterval changes the interval between successive acquisitions of metrics
// and their sending to the user interface. The default interval is one second.
func WithInterval(intv time.Duration) Option {
	return func(s *Server) {
		s.intv = intv
	}
}

// WithRoot changes the root path of the Statsviz user interface. The default
// root path is "/debug/statsviz".
func WithRoot(path string) Option {
	return func(s *Server) {
		s.root = path
	}
}

// WithTimeseriesPlot adds a new time series plot to Statsviz. This function can
// be called multiple times.
func WithTimeseriesPlot(tsp TimeSeriesPlot) Option {
	return func(s *Server) {
		s.userPlots = append(s.userPlots, plot.UserPlot{Scatter: tsp.timeseries})
	}
}

// Register registers the Statsviz HTTP handlers on the provided mux.
func (s *Server) Register(mux *http.ServeMux) {
	if s.plots == nil {
		// s is the zero value.
		s = NewServer()
	}

	mux.Handle(s.root+"/", s.Index())
	mux.HandleFunc(s.root+"/ws", s.Ws())
}

// intercept is a middleware that intercepts requests for plotsdef.js, which is
// generated dynamically based on the plots configuration. Other requests are
// forwarded as-is.
func intercept(h http.Handler, cfg *plot.Config) http.HandlerFunc {
	buf := &bytes.Buffer{}
	buf.WriteString("export default ")
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		panic("unexpected failure to encode plot definitions: " + err.Error())
	}
	buf.WriteString(";")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "js/plotsdef.js" {
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

// contentTypes forces the Content-Type HTTP header for certain files of some
// JavaScript libraries that have no extensions. Otherwise, the HTTP file server
// would serve them with "Content-Type: text/plain".
var contentTypes = map[string]string{
	"libs/js/popperjs-core2": "text/javascript",
	"libs/js/tippy.js@6":     "text/javascript",
}

// Index returns the index handler, which responds with the Statsviz user
// interface HTML page. By default, the handler is served at the path specified
// by the root. Use [WithRoot] to change the path.
func (s *Server) Index() http.HandlerFunc {
	prefix := strings.TrimSuffix(s.root, "/") + "/"
	assetsFS := http.FileServer(http.FS(static.Assets))
	return http.StripPrefix(prefix, intercept(assetsFS, s.plots.Config())).ServeHTTP
}

// Ws returns the WebSocket handler used by Statsviz to send application
// metrics. The underlying net.Conn is used to upgrade the HTTP server
// connection to the WebSocket protocol.
func (s *Server) Ws() http.HandlerFunc {
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

		// Ignore this error. This happens when the other end connection closes,
		// for example. We can't handle it in any meaningful way anyways.
		_ = s.sendStats(ws, s.intv)
	}
}

// sendStats sends runtime statistics over the WebSocket connection.
func (s *Server) sendStats(conn *websocket.Conn, frequency time.Duration) error {
	tick := time.NewTicker(frequency)
	defer tick.Stop()

	// If the WebSocket connection is initiated by an already open web UI
	// (started by a previous process, for example), then plotsdef.js won't be
	// requested. Call plots.Config() manually to ensure that s.plots internals
	// are correctly initialized.
	s.plots.Config()

	for range tick.C {
		w, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return err
		}
		if err := s.plots.WriteValues(w); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	}

	panic("unreachable")
}
