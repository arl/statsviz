package statsviz_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arl/statsviz"
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
	testIndex(t, statsviz.Index, "http://example.com/debug/statsviz/")
}

func TestIndexAtRoot(t *testing.T) {
	testIndex(t, statsviz.IndexAtRoot("/debug/"), "http://example.com/debug/")
	testIndex(t, statsviz.IndexAtRoot("/debug"), "http://example.com/debug/")
	testIndex(t, statsviz.IndexAtRoot("/"), "http://example.com/")
	testIndex(t, statsviz.IndexAtRoot("/test/"), "http://example.com/test/")
}
