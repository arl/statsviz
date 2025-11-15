package statsviz

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/arl/statsviz/internal/static"
)

func testIndex(t *testing.T, f http.Handler, url string) {
	t.Helper()

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)

	resp := w.Result()
	httpindex, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("http status %v, want %v", resp.StatusCode, http.StatusOK)
	}

	if resp.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("header[Content-Type] %s, want %s", resp.Header.Get("Content-Type"), "text/html; charset=utf-8")
	}

	fhtml, err := static.Assets().Open("index.html")
	if err != nil {
		t.Fatalf("couldn't read index.html from assets Fs: %v", err)
	}
	fsindex, _ := io.ReadAll(fhtml)

	if !bytes.Equal(fsindex, httpindex) {
		t.Errorf("read body is not that of index.html from assets")
	}
}

func newServer(tb testing.TB, opts ...Option) *Server {
	tb.Helper()

	srv, err := NewServer(opts...)
	if err != nil {
		tb.Fatal(err)
	}
	return srv
}

func TestIndex(t *testing.T) {
	t.Parallel()

	srv := newServer(t)
	testIndex(t, srv.Index(), "http://example.com/debug/statsviz/")
}

func TestRoot(t *testing.T) {
	t.Parallel()

	testIndex(t, newServer(t, Root("/debug/")).Index(), "http://example.com/debug/")
	testIndex(t, newServer(t, Root("/debug")).Index(), "http://example.com/debug/")
	testIndex(t, newServer(t, Root("/")).Index(), "http://example.com/")
	testIndex(t, newServer(t, Root("/test/")).Index(), "http://example.com/test/")
}

func testWs(t *testing.T, f http.Handler, URL string) {
	t.Helper()

	s := httptest.NewServer(f)
	defer s.Close()

	// Build a "ws://" url using the httptest server URL and the URL argument.
	u1, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	u2, err := url.Parse(URL)
	if err != nil {
		t.Fatal(err)
	}

	u1.Scheme = "ws"
	u1.Path = u2.Path

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u1.String(), nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// First message is the plots configuration.
	var cfg map[string]any
	if err := ws.ReadJSON(&cfg); err != nil {
		t.Fatalf("failed reading json from websocket: %v", err)
	}

	// Check the content of 2 consecutive payloads.
	for range 2 {
		// Verifies that we've received:
		// - 1 time series (cgo)
		// - 1 heatmap (sizeClasses).
		var msg struct {
			Event string `json:"event"`
			Data  struct {
				Series struct {
					CGo         []uint64 `json:"cgo"`
					SizeClasses []uint64 `json:"size-classes"`
				} `json:"series"`
			} `json:"data"`
		}

		if err := ws.ReadJSON(&msg); err != nil {
			t.Fatalf("failed reading json from websocket: %v", err)
		}

		// The time series must have one and only one element
		if len(msg.Data.Series.CGo) != 1 {
			t.Errorf("len(cgo) = %d, want 1", len(msg.Data.Series.CGo))
		}
		// Heatmaps should have many elements, check that there's more than one.
		if len(msg.Data.Series.SizeClasses) <= 1 {
			t.Errorf("len(sizeClasses) = %d, want > 1", len(msg.Data.Series.SizeClasses))
		}
	}
}

func TestWsCantUpgrade(t *testing.T) {
	url := "http://example.com/debug/statsviz/ws"

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	newServer(t).Ws()(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("responded %v to %q with non-websocket-upgradable conn, want %v", w.Result().StatusCode, url, http.StatusBadRequest)
	}
}

func testRegister(t *testing.T, f http.Handler, baseURL string) {
	testIndex(t, f, baseURL)
	ws := strings.TrimRight(baseURL, "/") + "/ws"
	testWs(t, f, ws)
}

func TestRegister(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		newServer(t).Register(mux)
		testRegister(t, mux, "http://example.com/debug/statsviz/")
	})

	t.Run("zero-value", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()

		var srv Server
		srv.Register(mux)
		testRegister(t, mux, "http://example.com/debug/statsviz/")
	})

	t.Run("root", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		newServer(t,
			Root(""),
		).Register(mux)

		testRegister(t, mux, "http://example.com/")
	})

	t.Run("root2", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		newServer(t,
			Root("/path/to/statsviz"),
		).Register(mux)

		testRegister(t, mux, "http://example.com/path/to/statsviz/")
	})

	t.Run("root+frequency", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		newServer(t,
			Root("/path/to/statsviz"),
			SendFrequency(100*time.Millisecond),
		).Register(mux)

		testRegister(t, mux, "http://example.com/path/to/statsviz/")
	})

	t.Run("non-positive frequency", func(t *testing.T) {
		t.Parallel()

		if _, err := NewServer(
			Root("/path/to/statsviz"),
			SendFrequency(-1),
		); err == nil {
			t.Errorf("NewServer() should have errored")
		}
	})
}

func TestRegisterDefault(t *testing.T) {
	mux := http.DefaultServeMux
	Register(mux)
	testRegister(t, mux, "http://example.com/debug/statsviz/")
}

func TestWsConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	t.Parallel()

	srv := newServer(t, SendFrequency(10*time.Millisecond))
	srv.Register(http.NewServeMux())

	s := httptest.NewServer(srv.Ws())
	defer s.Close()

	// Build a "ws://" url using the httptest server URL.
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Scheme = "ws"

	const numConns = 5
	const numMessages = 10

	errCh := make(chan error, numConns)

	for i := range numConns {
		go func(connID int) {
			ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				errCh <- err
				return
			}
			defer ws.Close()

			// First message is the plots configuration
			var cfg map[string]any
			if err := ws.ReadJSON(&cfg); err != nil {
				errCh <- err
				return
			}

			// Read multiple data messages to ensure state is being accessed
			for range numMessages {
				var msg map[string]any
				if err := ws.ReadJSON(&msg); err != nil {
					errCh <- err
					return
				}
			}

			errCh <- nil
		}(i)
	}

	for i := range numConns {
		if err := <-errCh; err != nil {
			t.Fatalf("connection %d failed: %v", i, err)
		}
	}
}
