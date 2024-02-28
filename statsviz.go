// Package statsviz allows visualizing Go runtime metrics data in real time in
// your browser.
//
// Register a Statsviz HTTP handlers with your server's [http.ServeMux]
// (preferred method):
//
//	mux := http.NewServeMux()
//	statsviz.Register(mux)
//
// Alternatively, you can register with [http.DefaultServeMux]:
//
//	ss := statsviz.Server{}
//	s.Register(http.DefaultServeMux)
//
// By default, Statsviz is served at http://host:port/debug/statsviz/. This, and
// other settings, can be changed by passing some [Option] to [NewServer].
//
// If your application is not already running an HTTP server, you need to start
// one. Add "net/http" and "log" to your imports, and use the following code in
// your main function:
//
//	go func() {
//	    log.Println(http.ListenAndServe("localhost:8080", nil))
//	}()
//
// Then open your browser and visit http://localhost:8080/debug/statsviz/.
//
// # Advanced usage:
//
// If you want more control over Statsviz HTTP handlers, examples are:
//   - you're using some HTTP framework
//   - you want to place Statsviz handler behind some middleware
//
// then use [NewServer] to obtain a [Server] instance. Both the [Server.Index] and
// [Server.Metrics]() methods return [http.HandlerFunc].
//
//	srv, err := statsviz.NewServer(); // Create server or handle error
//	srv.Index()                       // UI (dashboard) http.HandlerFunc
//	srv.Metrics()                     // Metrics http.HandlerFunc
package statsviz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arl/statsviz/internal/plot"
	"github.com/arl/statsviz/internal/static"
)

const (
	defaultRoot         = "/debug/statsviz"
	defaultMetrics      = "metrics"
	defaultSendInterval = time.Second
)

// RegisterDefault registers the Statsviz HTTP handlers on [http.DefaultServeMux].
//
// RegisterDefault should not be used in production.
func RegisterDefault(opts ...Option) error {
	return Register(http.DefaultServeMux)
}

// Register registers the Statsviz HTTP handlers on the provided mux.
func Register(mux *http.ServeMux, opts ...Option) error {
	srv, err := NewServer(opts...)
	if err != nil {
		return err
	}
	srv.Register(mux)
	return nil
}

// Server is the core component of Statsviz. It collects and periodically
// updates metrics data and provides two essential HTTP handlers:
//   - the Index handler serves Statsviz user interface, allowing you to
//     visualize runtime metrics on your browser.
//   - The Metrics handler establishes a data connection allowing the connected
//     browser to receive metrics updates from the server.
//
// The zero value is not a valid Server, use NewServer to create a valid one.
type Server struct {
	intv      time.Duration // interval between consecutive metrics emission
	root      string        // HTTP path root
	metrics   string        // http path for metrics
	plots     *plot.List    // plots shown on the user interface
	userPlots []plot.UserPlot
}

// NewServer constructs a new Statsviz Server with the provided options, or the
// default settings.
//
// Note that once the server is created, its HTTP handlers needs to be registered
// with some HTTP server. You can either use the Register method or register yourself
// the Index and Ws handlers.
func NewServer(opts ...Option) (*Server, error) {
	s := &Server{
		intv:    defaultSendInterval,
		root:    defaultRoot,
		metrics: defaultMetrics,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	pl, err := plot.NewList(s.userPlots)
	if err != nil {
		return nil, err
	}
	s.plots = pl
	return s, nil
}

// Option is a configuration option for the Server.
type Option func(*Server) error

// SendFrequency changes the interval between successive acquisitions of metrics
// and their sending to the user interface. The default interval is one second.
func SendFrequency(intv time.Duration) Option {
	return func(s *Server) error {
		if intv <= 0 {
			return fmt.Errorf("frequency must be a positive integer")
		}
		s.intv = intv
		return nil
	}
}

// Root changes the root path of the Statsviz user interface.
// The default is "/debug/statsviz".
func Root(path string) Option {
	return func(s *Server) error {
		s.root = path
		return nil
	}
}

// MetricsPath changes the metrics path of the Statsviz user interface.
// The default is root+"/metrics".
func MetricsPath(path string) Option {
	return func(s *Server) error {
		s.metrics = path
		return nil
	}
}

// TimeseriesPlot adds a new time series plot to Statsviz. This options can
// be added multiple times.
func TimeseriesPlot(tsp TimeSeriesPlot) Option {
	return func(s *Server) error {
		s.userPlots = append(s.userPlots, plot.UserPlot{Scatter: tsp.timeseries})
		return nil
	}
}

// Register registers the Statsviz HTTP handlers on the provided mux.
func (s *Server) Register(mux *http.ServeMux) {
	mux.Handle(s.root+"/", s.Index())
	if s.metrics == "" {
		s.metrics = defaultMetrics
	}
	mux.HandleFunc(s.root+"/"+s.metrics, s.Metrics())
}

// intercept is a middleware that intercepts requests for plotsdef.js, which is
// generated dynamically based on the plots configuration. Other requests are
// forwarded as-is.
func intercept(h http.Handler, cfg *plot.Config, extraConfig map[string]any) http.HandlerFunc {
	var plotsdefjs []byte
	//Using parentheses helps gc
	{
		buf := bytes.Buffer{}
		buf.WriteString("export default ")
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		var encodeValue any = cfg
		if len(extraConfig) > 0 {
			encodeValue1 := map[string]any{
				"series": cfg.Series,
				"events": cfg.Events,
			}
			for k, v := range extraConfig {
				encodeValue1[k] = v
			}
			encodeValue = encodeValue1
		}
		if err := enc.Encode(encodeValue); err != nil {
			panic("unexpected failure to encode plot definitions: " + err.Error())
		}
		buf.WriteString(";")
		plotsdefjs = buf.Bytes()
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "js/plotsdef.js" {
			w.Header().Add("Content-Length", strconv.Itoa(len(plotsdefjs)))
			w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
			w.Write(plotsdefjs)
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

// Returns an FS serving the embedded assets, or the assets directory if
// STATSVIZ_DEBUG contains the 'asssets' key.
func assetsFS() http.FileSystem {
	assets := http.FS(static.Assets)

	vdbg := os.Getenv("STATSVIZ_DEBUG")
	if vdbg == "" {
		return assets
	}

	kvs := strings.Split(vdbg, ";")
	for _, kv := range kvs {
		k, v, found := strings.Cut(strings.TrimSpace(kv), "=")
		if !found {
			panic("invalid STATSVIZ_DEBUG value: " + kv)
		}
		if k == "assets" {
			dir := filepath.Join(v)
			return http.Dir(dir)
		}
	}

	return assets
}

// Index returns the index handler, which responds with the Statsviz user
// interface HTML page. By default, the handler is served at the path specified
// by the root. The default path is "/debug/statsviz/". Use [Root] to change the path.
func (s *Server) Index() http.HandlerFunc {
	prefix := strings.TrimSuffix(s.root, "/") + "/"
	assets := http.FileServer(assetsFS())
	// defer initialization until the actual request, so that the Server's properties(s.xxx) are fixed
	once := sync.Once{}
	var realHandler http.HandlerFunc
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		once.Do(func() {
			realHandler = intercept(assets, s.plots.Config(), map[string]any{
				"sendFrequency": s.intv.Milliseconds(),
				"metricsPath":   s.metrics,
			})
		})
		// the sse protocol in github.com/soheilhy/cmux and other frameworks may reuse other requests,
		// actively close to avoid bugs.
		writer.Header().Add("Connection", "close")
		realHandler.ServeHTTP(writer, request)
	})
	return http.StripPrefix(prefix, handler).ServeHTTP
}

// Ws returns the long connection handler used by Statsviz to send application metrics.
// The default path is root+"/ws". Use [MetricsPath] to change the path.
// Deprecated: use Metrics instead
func (s *Server) Ws() http.HandlerFunc {
	println("statsviz.Server.Ws() is deprecated, use statsviz.Server.Metrics() instead")
	//if you use the websockt version of writing, we will change the default path to ws
	if s.metrics == "" || s.metrics == defaultMetrics {
		s.metrics = "ws"
	}
	return s.Metrics()
}

// Metrics returns the long connection handler used by Statsviz to send application metrics.
// The default path is root+"/metrics". Use [MetricsPath] to change the path.
func (s *Server) Metrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept"), "/event-stream") {
			// If the connection is initiated by an already open web UI
			// (started by a previous process, for example), then plotsdef.js won't be
			// requested. Call plots.Config() manually to ensure that s.plots internals
			// are correctly initialized.
			s.plots.Config()

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			s.startTransfer(w)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This endpoint only supports text/event-stream requests"))
	}
}

func (s *Server) startTransfer(w io.Writer) {
	buffer := bytes.Buffer{}
	buffer.WriteString("data: ")
	callData := func() error {
		if err := s.plots.WriteValues(&buffer); err == nil {
			_, err = w.Write(buffer.Bytes())
			if err != nil {
				return err
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		} else {
			return err
		}
		return nil
	}
	//the first time it was sent immediately
	err := callData()
	if err != nil {
		return
	}
	tick := time.NewTicker(s.intv)
	defer tick.Stop()
	for range tick.C {
		if callData() != nil {
			return
		}
	}
}
