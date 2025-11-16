//go:build go1.25
// +build go1.25

package statsviz

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"testing/synctest"

	"github.com/gorilla/websocket"
)

func TestWsConcurrent(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		srv := newServer(t)
		srv.Register(http.NewServeMux())

		li := fakeNetListen()
		s := &httptest.Server{
			Listener: li,
			Config:   &http.Server{Handler: srv.Ws()},
		}
		s.Start()

		// Build a "ws://" url using the httptest server URL.
		u, err := url.Parse(s.URL)
		if err != nil {
			t.Fatal(err)
		}
		u.Scheme = "ws"

		const numConns = 10
		const numMessages = 200

		errCh := make(chan error, numConns)

		synctestWsDialer := websocket.Dialer{
			NetDialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				c := li.connect()
				return c, nil
			},
		}

		for i := range numConns {
			go func(connID int) {
				ws, _, err := synctestWsDialer.Dial(u.String(), nil)
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

		srv.Close()
		s.Close()

		synctest.Wait()
	})
}
