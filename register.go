package statsviz

import (
	"fmt"
	"net/http"
	"time"
)

// RegisterDefault registers statsviz HTTP handlers on the default serve mux.
//
// Note this is not advised on a production server, unless it only serves on
// localhost.
func RegisterDefault(opts ...OptionFunc) error {
	return Register(http.DefaultServeMux, opts...)
}

// Root sets the root path of statsviz handlers.
func Root(root string) OptionFunc {
	return func(s *server) error {
		s.root = root
		return nil
	}
}

// SendFrequency defines the frequency at which statistics are sent from the
// application to the HTML page.
func SendFrequency(freq time.Duration) OptionFunc {
	return func(s *server) error {
		if freq <= 0 {
			return fmt.Errorf("frequency must be a positive integer")
		}
		s.freq = freq
		return nil
	}
}

// An OptionFunc is a server configuration option.
type OptionFunc func(s *server) error

const (
	defaultRoot          = "/debug/statsviz"
	defaultSendFrequency = time.Second
)

// Register registers statsviz HTTP handlers on the provided mux.
func Register(mux *http.ServeMux, opts ...OptionFunc) error {
	s := &server{
		mux:  mux,
		root: defaultRoot,
		freq: defaultSendFrequency,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}

	s.register()
	return nil
}

type server struct {
	mux  *http.ServeMux
	freq time.Duration
	root string
}

func (s *server) register() {
	s.mux.Handle(s.root+"/", IndexAtRoot(s.root))
	s.mux.HandleFunc(s.root+"/ws", NewWsHandler(s.freq))
}
