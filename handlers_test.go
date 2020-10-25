package statsviz

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arl/statsviz/websocket"
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
	testIndex(t, Index, "http://example.com/debug/statsviz/")
}

func TestIndexAtRoot(t *testing.T) {
	testIndex(t, IndexAtRoot("/debug/"), "http://example.com/debug/")
	testIndex(t, IndexAtRoot("/debug"), "http://example.com/debug/")
	testIndex(t, IndexAtRoot("/"), "http://example.com/")
	testIndex(t, IndexAtRoot("/test/"), "http://example.com/test/")
}

func testWs(t *testing.T, f http.Handler, url string) {
	t.Helper()

	s := httptest.NewServer(f)
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
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
	testWs(t, http.HandlerFunc(Ws), "http://example.com/debug/statsviz/ws")
}
