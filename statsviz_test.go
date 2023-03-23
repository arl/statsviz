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

func TestIndex(t *testing.T) {
	t.Parallel()

	se := NewEndpoint()
	testIndex(t, se.Index(), "http://example.com/debug/statsviz/")
}

func TestRoot(t *testing.T) {
	t.Parallel()

	testIndex(t, NewEndpoint(WithRoot("/debug/")).Index(), "http://example.com/debug/")
	testIndex(t, NewEndpoint(WithRoot("/debug")).Index(), "http://example.com/debug/")
	testIndex(t, NewEndpoint(WithRoot("/")).Index(), "http://example.com/")
	testIndex(t, NewEndpoint(WithRoot("/test/")).Index(), "http://example.com/test/")
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

	// Check the content of 2 consecutive payloads.
	for i := 0; i < 2; i++ {

		// Verifies that we've received 1 time series (goroutines) and one
		// heatmap (sizeClasses).
		var data struct {
			Goroutines  []uint64 `json:"goroutines"`
			SizeClasses []uint64 `json:"size-classes"`
		}
		if err := ws.ReadJSON(&data); err != nil {
			t.Fatalf("failed reading json from websocket: %v", err)
		}

		// The time series must have one and only one element
		if len(data.Goroutines) != 1 {
			t.Errorf("len(goroutines) = %d, want 1", len(data.Goroutines))
		}
		// Heatmaps should have many elements, check that there's more than one.
		if len(data.SizeClasses) <= 1 {
			t.Errorf("len(sizeClasses) = %d, want > 1", len(data.SizeClasses))
		}
	}
}

func TestWs(t *testing.T) {
	t.Parallel()

	testWs(t, NewEndpoint().Ws(), "http://example.com/debug/statsviz/ws")
}

func TestWsCantUpgrade(t *testing.T) {
	url := "http://example.com/debug/statsviz/ws"

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	NewEndpoint().Ws()(w, req)

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
		NewEndpoint().Register(mux)
		testRegister(t, mux, "http://example.com/debug/statsviz/")
	})

	t.Run("root", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		NewEndpoint(WithRoot("")).Register(mux)

		testRegister(t, mux, "http://example.com/")
	})

	t.Run("root2", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		NewEndpoint(WithRoot("/path/to/statsviz")).Register(mux)

		testRegister(t, mux, "http://example.com/path/to/statsviz/")
	})

	t.Run("root+frequency", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		NewEndpoint(
			WithRoot("/root/to/statsviz"),
			WithInterval(100*time.Millisecond),
		).Register(mux)

		testRegister(t, mux, "http://example.com/root/to/statsviz/")
	})

	t.Run("non-positive frequency", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		NewEndpoint(
			WithRoot("/root/to/statsviz"),
			WithInterval(0),
		).Register(mux)
	})
}

func TestRegisterDefault(t *testing.T) {
	NewEndpoint().Register(http.DefaultServeMux)
	testRegister(t, http.DefaultServeMux, "http://example.com/debug/statsviz/")
}

func Test_intercept(t *testing.T) {
	// Check that the file server has been 'hijacked'.
	// 'plotsdef.js' is generated at runtime, it doesn't actually exist, it is generated on the fly.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/debug/statsviz/js/plotsdef.js", nil)

	se := NewEndpoint()
	intercept(se.Index(), se.plots.Config())(w, req)

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
	e := NewEndpoint(WithRoot(root))
	httpfs := e.Index()

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

		ct := w.HeaderMap.Get("Content-Type")
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
