package statsviz

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func testIndex(t *testing.T, f http.Handler, url string) {
	t.Helper()

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("http status %v, want %v", resp.StatusCode, http.StatusOK)
	}

	if resp.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("header[Content-Type] %s, want %s", resp.Header.Get("Content-Type"), "text/html; charset=utf-8")
	}

	if !strings.Contains(string(body), "goroutines") {
		t.Errorf("body doesn't contain %q", "goroutines")
	}
}

func TestIndex(t *testing.T) {
	t.Parallel()

	testIndex(t, Index, "http://example.com/debug/statsviz/")
}

func TestIndexAtRoot(t *testing.T) {
	t.Parallel()

	testIndex(t, IndexAtRoot("/debug/"), "http://example.com/debug/")
	testIndex(t, IndexAtRoot("/debug"), "http://example.com/debug/")
	testIndex(t, IndexAtRoot("/"), "http://example.com/")
	testIndex(t, IndexAtRoot("/test/"), "http://example.com/test/")
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

	// Wait for 2 messages and check that the payload is what we expect.
	for i := 0; i < 2; i++ {
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		var stats stats
		if err := json.Unmarshal(p, &stats); err != nil {
			t.Fatal(err)
		}
	}
}

func TestWs(t *testing.T) {
	t.Parallel()

	testWs(t, http.HandlerFunc(Ws), "http://example.com/debug/statsviz/ws")
}

func TestWsCantUpgrade(t *testing.T) {
	url := "http://example.com/debug/statsviz/ws"

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	Ws(w, req)

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
		if err := Register(mux); err != nil {
			t.Fatal(err)
		}

		testRegister(t, mux, "http://example.com/debug/statsviz/")
	})

	t.Run("root", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		if err := Register(mux, Root("")); err != nil {
			t.Fatal(err)
		}

		testRegister(t, mux, "http://example.com/")
	})

	t.Run("root2", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		if err := Register(mux, Root("/root/to/statsviz")); err != nil {
			t.Fatal(err)
		}

		testRegister(t, mux, "http://example.com/root/to/statsviz/")
	})

	t.Run("root+frequency", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		err := Register(mux, Root("/root/to/statsviz"), SendFrequency(100*time.Millisecond))
		if err != nil {
			t.Fatal(err)
		}

		testRegister(t, mux, "http://example.com/root/to/statsviz/")
	})

	t.Run("non-positive frequency", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		err := Register(mux, Root("/root/to/statsviz"), SendFrequency(0))
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestRegisterDefault(t *testing.T) {
	if err := RegisterDefault(); err != nil {
		t.Fatal(err)
	}

	testRegister(t, http.DefaultServeMux, "http://example.com/debug/statsviz/")
}
