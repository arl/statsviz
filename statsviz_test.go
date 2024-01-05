package statsviz

import (
	"bytes"
	"io"
	"io/fs"
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
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("http status %v, want %v", resp.StatusCode, http.StatusOK)
	}

	if resp.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("header[Content-Type] %s, want %s", resp.Header.Get("Content-Type"), "text/html; charset=utf-8")
	}

	html, err := static.Assets.ReadFile("index.html")
	if err != nil {
		t.Fatalf("couldn't read index.html from assets Fs: %v", err)
	}

	if !bytes.Equal(html, body) {
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

func testWs(t *testing.T, s *httptest.Server, URL string, number int, checkData func() any, check func(*testing.T, any)) {
	t.Helper()

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

	// Check the content of 2 consecutive payloads.
	for i := 0; i < number; i++ {
		data := checkData()
		if err := ws.ReadJSON(data); err != nil {
			t.Fatalf("failed reading json from websocket: %v", err)
			return
		}
		check(t, data)
	}
}

func TestWs(t *testing.T) {
	t.Parallel()
	// Verifies that we've received 1 time series (goroutines) and one
	// heatmap (sizeClasses).
	type dataType struct {
		Goroutines  []uint64 `json:"goroutines"`
		SizeClasses []uint64 `json:"size-classes"`
	}

	s := httptest.NewServer(newServer(t).Ws())
	defer s.Close()
	testWs(t, s, "http://example.com/debug/statsviz/ws", 2, func() any {
		return &dataType{}
	}, func(t *testing.T, data1 any) {
		data := data1.(*dataType)
		// The time series must have one and only one element
		if len(data.Goroutines) != 1 {
			t.Errorf("len(goroutines) = %d, want 1", len(data.Goroutines))
		}
		// Heatmaps should have many elements, check that there's more than one.
		if len(data.SizeClasses) <= 1 {
			t.Errorf("len(sizeClasses) = %d, want > 1", len(data.SizeClasses))
		}
	})
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
	type dataType struct {
		Goroutines  []uint64 `json:"goroutines"`
		SizeClasses []uint64 `json:"size-classes"`
	}
	s := httptest.NewServer(f)
	defer s.Close()
	testWs(t, s, ws, 2, func() any {
		return &dataType{}
	}, func(t *testing.T, data1 any) {
		data := data1.(*dataType)
		if len(data.Goroutines) != 1 {
			t.Errorf("len(goroutines) = %d, want 1", len(data.Goroutines))
		}
		if len(data.SizeClasses) <= 1 {
			t.Errorf("len(sizeClasses) = %d, want > 1", len(data.SizeClasses))
		}
	})
}

func TestRegister(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		newServer(t).Register(mux)
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
	type Data struct {
		Test []float64 `json:"test"`
	}

	makeTestPlot := func() (TimeSeriesPlot, *int) {
		num := 0
		build, _ := TimeSeriesPlotConfig{
			Title: "test",
			Name:  "test",
			Series: []TimeSeries{
				{
					Name: "1",
					GetValue: func() float64 {
						num++
						return 1
					},
				},
			},
		}.Build()
		return build, &num
	}

	t.Run("customizePlot", func(t *testing.T) {
		t.Parallel()
		plot, i := makeTestPlot()
		s := httptest.NewServer(newServer(t,
			TimeseriesPlot(plot),
		).Ws())
		defer s.Close()
		go testWs(t, s, "http://example.com/debug/statsviz/ws", 2, func() any { return &Data{} }, func(t *testing.T, a any) {})
		testWs(t, s, "http://example.com/debug/statsviz/ws", 2, func() any {
			return &Data{}
		}, func(t *testing.T, a any) {
			data := a.(*Data)
			if len(data.Test) != 1 {
				t.Fatalf("customizePlot failed,call num %d expect 1", len(data.Test))
			}
		})
		time.Sleep(100 * time.Millisecond)
		if *i != 4 {
			t.Fatalf("dataCache2 failed,call num %d expect 4", *i)
		}
	})

	t.Run("dataCache", func(t *testing.T) {
		t.Parallel()
		plot, i := makeTestPlot()
		s := httptest.NewServer(newServer(t,
			EnableDataCache(),
			TimeseriesPlot(plot),
		).Ws())
		defer s.Close()
		go testWs(t, s, "http://example.com/debug/statsviz/ws", 1, func() any { return &Data{} }, func(t *testing.T, a any) {})
		go testWs(t, s, "http://example.com/debug/statsviz/ws", 2, func() any { return &Data{} }, func(t *testing.T, a any) {})
		go testWs(t, s, "http://example.com/debug/statsviz/ws", 3, func() any { return &Data{} }, func(t *testing.T, a any) {})
		go testWs(t, s, "http://example.com/debug/statsviz/ws", 4, func() any { return &Data{} }, func(t *testing.T, a any) {})
		testWs(t, s, "http://example.com/debug/statsviz/ws", 4, func() any { return &Data{} }, func(t *testing.T, a any) {})
		time.Sleep(100 * time.Millisecond)
		if *i != 4 {
			t.Fatalf("dataCache failed,call num %d expect 4", *i)
		}
	})

	t.Run("dataCache2", func(t *testing.T) {
		t.Parallel()
		plot, i := makeTestPlot()
		s := httptest.NewServer(newServer(t,
			EnableDataCache(),
			TimeseriesPlot(plot),
		).Ws())
		defer s.Close()
		go testWs(t, s, "http://example.com/debug/statsviz/ws", 2, func() any { return &Data{} }, func(t *testing.T, a any) {})
		testWs(t, s, "http://example.com/debug/statsviz/ws", 2, func() any { return &Data{} }, func(t *testing.T, a any) {})
		time.Sleep(100 * time.Millisecond)
		if *i != 2 {
			t.Fatalf("dataCache failed,call num %d expect 2", *i)
		}
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

func Test_intercept(t *testing.T) {
	// Check that the file server has been 'hijacked'.
	// 'plotsdef.js' is generated at runtime, it doesn't actually exist, it is generated on the fly.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/debug/statsviz/js/plotsdef.js", nil)

	srv := newServer(t)
	intercept(srv.Index(), srv.plots.Config())(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("http status %v, want %v", resp.StatusCode, http.StatusOK)
	}

	contentType := "text/javascript; charset=utf-8"
	if resp.Header.Get("Content-Type") != contentType {
		t.Errorf("header[Content-Type] %s, want %s", resp.Header.Get("Content-Type"), contentType)
	}
}

func TestContentTypeIsSet(t *testing.T) {
	// Check that "Content-Type" headers on the assets we serve are all set to
	// something more specific than "text/plain" because that'd make the page be
	// rejected in certain 'strict' environments.
	const root = "/some/root/path"
	srv := newServer(t, Root(root))
	httpfs := srv.Index()

	requested := []string{}

	// While we walk the embedded assets filesystem, control the header on the
	// http filesystem server.
	_ = fs.WalkDir(static.Assets, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || path == "fs.go" || path == "index.html" {
			return nil
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, root+"/"+path, nil)

		httpfs(w, r)
		res := w.Result()
		if res.StatusCode != 200 && path != "index.html" {
			t.Errorf("GET %q returned HTTP %d, want 200", path, res.StatusCode)
			return nil
		}

		ct := res.Header.Get("Content-Type")
		if ct == "" || strings.Contains(ct, "text/plain") {
			t.Errorf(`GET %q has incorrect header "Content-Type = %s"`, path, ct)
			return nil
		}

		if testing.Verbose() {
			t.Logf("%q Content-Type %q", path, ct)
		}
		requested = append(requested, path)
		return nil
	})

	// Verify that all files in contentTypes map have been requested. This is to
	// keep the map aligned with the actual content of the static/ dir.
	for path := range contentTypes {
		found := false
		for i := range requested {
			if requested[i] == path {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("contentTypes[%v] matches no files in the static/ dir", path)
		}
	}
}
